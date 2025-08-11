package db

import (
	"context"

	"github.com/Gilf4/grpcChat/auth/internal/domain/models"
)

type UserStorage struct {
}

func (s *UserStorage) Create(ctx context.Context, email string, passHash []byte, name string) (id int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func New() *UserStorage {
	return &UserStorage{}
}
