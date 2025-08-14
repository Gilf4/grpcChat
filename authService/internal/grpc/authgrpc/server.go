package authgrpc

import (
	"context"
	"errors"

	"github.com/Gilf4/grpcChat/auth/internal/repository/db"
	"github.com/Gilf4/grpcChat/auth/internal/services/auth"
	authv1 "github.com/Gilf4/grpcChat/protos/gen/go/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	Register(ctx context.Context, email, password, name string) (int64, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	email := req.GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetPassword()
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	accessToken, refreshToken, err := s.auth.Login(ctx, email, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &authv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	email := req.GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	password := req.GetPassword()
	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	name := req.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	userId, err := s.auth.Register(ctx, email, password, name)
	if err != nil {
		if errors.Is(err, db.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &authv1.RegisterResponse{UserId: userId}, nil
}

func (s *serverAPI) RefreshAccessToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	refreshToken := req.GetRefreshToken()
	if refreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	accessToken, err := s.auth.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		}
		return nil, status.Error(codes.Internal, "failed to refresh access token")
	}

	return &authv1.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}
