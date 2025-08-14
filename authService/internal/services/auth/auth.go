package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Gilf4/grpcChat/auth/internal/domain/models"
	"github.com/Gilf4/grpcChat/auth/internal/lib/jwt"
	"github.com/Gilf4/grpcChat/auth/internal/lib/refresh"
	"github.com/Gilf4/grpcChat/auth/internal/repository/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

type UserRepository interface {
	Create(ctx context.Context, email string, passHash []byte, name string) (int64, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByID(ctx context.Context, id int64) (models.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, userID int64, refreshToken string, ttl time.Duration) error
	GetByToken(ctx context.Context, token string) (*models.Session, error)
	Delete(ctx context.Context, token string) error
}

type Auth struct {
	log             *slog.Logger
	userRepo        UserRepository
	sessionRepo     SessionRepository
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	jwtSecret       string
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func New(
	log *slog.Logger,
	userRepo UserRepository,
	sessionRepo SessionRepository,
	AccessTokenTTL time.Duration,
	RefreshTokenTTL time.Duration,
	jwtSecret string,
) *Auth {
	return &Auth{
		log:             log,
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		accessTokenTTL:  AccessTokenTTL,
		refreshTokenTTL: RefreshTokenTTL,
		jwtSecret:       jwtSecret,
	}
}

func (a *Auth) Login(ctx context.Context, email, pass string) (string, string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email))

	log.Info("attempting to login user")

	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, db.ErrUserNotFound) {
			log.Warn("user not found", "error", err.Error())

			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", "error", err.Error())

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := VerifyPassword(user.PassHash, pass); err != nil {
		log.Info("invalid credentials", "error", err.Error())

		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	accessToken, err := jwt.NewToken(&user, a.accessTokenTTL, a.jwtSecret)
	if err != nil {
		log.Error("failed to generate access token", "error", err.Error())

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refreshToken, err := refresh.GenerateToken()
	if err != nil {
		log.Error("failed to generate refresh token", "error", err.Error())

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = a.sessionRepo.Create(ctx, user.ID, refreshToken, a.refreshTokenTTL)
	if err != nil {
		log.Error("failed to create refresh token in db", "error", err.Error())

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
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

	id, err := a.userRepo.Create(ctx, email, passHash, name)
	if err != nil {
		log.Error("failed to save user", "error", err.Error())

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	const op = "Auth.RefreshAccessToken"

	log := a.log.With(
		slog.String("op", op),
	)

	session, err := a.sessionRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, db.ErrTokenNotFound) {
			log.Warn("refresh token not found")
			return "", ErrInvalidRefreshToken
		}
		log.Error("failed to get session", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if time.Now().After(session.ExpiresAt) {
		log.Warn("refresh token expired")
		return "", ErrRefreshTokenExpired
	}

	user, err := a.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		log.Error("failed to get user", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := jwt.NewToken(&user, a.accessTokenTTL, a.jwtSecret)
	if err != nil {
		log.Error("failed to create access token", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("access token refreshed successfully")

	return accessToken, nil
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
