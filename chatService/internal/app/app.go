package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/Gilf4/grpcChat/chat/internal/app/grpc"
	"github.com/Gilf4/grpcChat/chat/internal/config"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	grpcApp := grpcapp.New(log, cfg.GRPC.Port)

	return &App{
		GRPCServer: grpcApp,
	}
}
