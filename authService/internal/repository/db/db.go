package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Gilf4/grpcChat/auth/internal/config"
	"github.com/Gilf4/grpcChat/auth/internal/domain/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	db *pgxpool.Pool
}

func NewUserRepository(ctx context.Context, dbCfg *config.DBConfig) (*UserStorage, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &UserStorage{db: pool}, nil
}

func (s *UserStorage) Create(
	ctx context.Context,
	email string,
	passHash []byte,
	name string,
) (int64, error) {
	op := "repo.db.Create"

	query := `
		INSERT INTO users (email, pass_hash, name)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int64

	err := s.db.QueryRow(ctx, query, email, passHash, name).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (models.User, error) {
	op := "repo.db.GetByEmail"

	query := `
		SELECT id, email, pass_hash, name
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
		&user.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
