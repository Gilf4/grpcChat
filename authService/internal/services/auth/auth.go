package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Gilf4/grpcChat/auth/internal/lib/jwt"
	"github.com/Gilf4/grpcChat/auth/internal/repository/db"
	"github.com/Gilf4/grpcChat/auth/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email string, passHash []byte, name string) (id int64, err error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type Auth struct {
	log       *slog.Logger
	repo      UserRepository
	tokenTTL  time.Duration
	jwtSecret string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(log *slog.Logger, repo UserRepository, tokenTTL time.Duration, jwtSecret string) *Auth {
	return &Auth{
		log:       log,
		repo:      repo,
		tokenTTL:  tokenTTL,
		jwtSecret: jwtSecret,
	}
}

func (a *Auth) Login(ctx context.Context, email, pass string) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email))

	log.Info("attempting to login user")

	user, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			a.log.Warn("user not found", "error", err.Error())

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", "error", err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := VerifyPassword(user.PasswordHash, pass); err != nil {
		a.log.Info("invalid credentials", "error", err.Error())

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, a.tokenTTL, a.jwtSecret)
	if err != nil {
		a.log.Error("failed to generate token", "error", err.Error())

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email, pass, name string) (int64, error) {
	const op = "Auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email))

	log.Info("attempting to register user")

	passHash, err := HashPassword(pass)
	if err != nil {
		log.Error("failed to generate password hash", "error", err.Error())

		return 0, err
	}

	id, err := a.repo.CreateUser(ctx, email, passHash, name)
	if err != nil {
		log.Error("failed to save user", "error", err.Error())

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
