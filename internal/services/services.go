package services

import (
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/services/users"
	"github.com/Tel3scop/auth/internal/storages"
)

// Container структура контейнера
type Container struct {
	Users users.Service
}

// New создать новый сервис
func New(cfg *config.Config, storages *storages.Storage) *Container {
	return &Container{Users: users.NewService(cfg, storages)}
}
