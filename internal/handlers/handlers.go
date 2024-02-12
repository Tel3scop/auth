package handlers

import (
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run запуск всех хендлеров
func Run() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	userAPI.RegisterUserV1Server(s, &userServer{})

	return s
}
