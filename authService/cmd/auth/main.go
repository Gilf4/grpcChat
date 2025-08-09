package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gilf4/grpcChat/auth/internal/app"
	"github.com/Gilf4/grpcChat/auth/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	// TODO: init database

	log.Info(
		"starting application:",
		slog.Any("", cfg),
	)

	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTL, cfg.JWTSecret)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signal := <-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped", "signal", signal)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
