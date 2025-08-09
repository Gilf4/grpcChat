package db

import (
	"context"
	"github.com/Gilf4/grpcChat/auth/internal/repository/models"
)

type Storage struct {
}

func (s Storage) CreateUser(ctx context.Context, email string, passHash []byte, name string) (id int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (s Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func New() *Storage {
	return &Storage{}
}
