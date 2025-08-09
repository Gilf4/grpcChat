package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Gilf4/grpcChat/auth/internal/app/grpcApp"
	"github.com/Gilf4/grpcChat/auth/internal/repository/db"
	"github.com/Gilf4/grpcChat/auth/internal/services/auth"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
	JWTSecret string,
) *App {
	userRepository := db.New()
	authService := auth.New(log, userRepository, tokenTTL, JWTSecret)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCServer: grpcApp,
	}
}
