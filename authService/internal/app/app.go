package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/Gilf4/grpcChat/auth/internal/app/grpcApp"
	"github.com/Gilf4/grpcChat/auth/internal/config"
	"github.com/Gilf4/grpcChat/auth/internal/repository/db"
	"github.com/Gilf4/grpcChat/auth/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	userRepository, err := db.NewUserRepository(ctx, &cfg.DB)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, userRepository, cfg.TokenTTL, cfg.JWTSecret)

	grpcApp := grpcapp.New(log, cfg.GRPC.Port, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
