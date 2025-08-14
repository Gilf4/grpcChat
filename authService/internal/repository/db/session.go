package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Gilf4/grpcChat/auth/internal/config"
	"github.com/Gilf4/grpcChat/auth/internal/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStorage struct {
	db *pgxpool.Pool
}

func NewSessionRepository(ctx context.Context, dbCfg *config.DBConfig) (*SessionStorage, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &SessionStorage{db: pool}, nil
}

func (s *SessionStorage) Create(ctx context.Context, userID int64, refreshToken string, ttl time.Duration) error {
	op := "repo.Session.Create"

	query := `
		INSERT INTO sessions (user_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := s.db.Exec(ctx, query, userID, refreshToken, time.Now().Add(ttl))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *SessionStorage) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	op := "repo.Sessions.GetByToken"

	query := `
		SELECT id, user_id, refresh_token, expires_at
		FROM sessions
		WHERE refresh_token = $1
	`

	var session models.Session
	err := s.db.QueryRow(ctx, query, token).Scan(&session.ID, &session.UserID, &session.RefreshToken, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Session{}, fmt.Errorf("%s: %w", op, ErrTokenNotFound)
		}
		return &models.Session{}, fmt.Errorf("%s: %w", op, err)
	}
	return &session, nil
}

func (s *SessionStorage) Delete(ctx context.Context, token string) error {
	op := "repo.Sessions.Delete"

	query := `
		DELETE FROM sessions
		WHERE refresh_token = $1
	`
	_, err := s.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}
