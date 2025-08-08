package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Gilf4/grpcChat/auth/internal/app/grpcApp"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, tokenTTL time.Duration) *App {
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
