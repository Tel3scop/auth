package handlers

import (
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/services"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Run запуск всех хендлеров
func Run(cfg *config.Config, serviceContainer *services.Container) *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	userAPI.RegisterUserV1Server(s, &userServer{
		services: serviceContainer,
		cfg:      cfg,
	})

	return s

}
