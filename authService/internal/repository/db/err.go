package db

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrTokenNotFound = errors.New("refreshToken not found")
)
