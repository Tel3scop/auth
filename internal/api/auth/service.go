package auth

import (
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/pkg/auth_v1"
)

// Implementation структура для работы с хэндерами авторизации
type Implementation struct {
	auth_v1.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation новый экземпляр структуры Implementation
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
