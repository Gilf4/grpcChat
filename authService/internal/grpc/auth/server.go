package auth

import (
	"context"

	"github.com/Gilf4/grpcChat/protos/gen/go/auth/v1"
	"google.golang.org/grpc"
)

type serverAPI struct {
	auth.UnimplementedAuthServer
}

func Register(gRPCServer *grpc.Server) {
	auth.RegisterAuthServer(gRPCServer, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	panic("implement me")
}
