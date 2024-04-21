package user

import (
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/pkg/user_v1"
)

// Implementation структура для работы с хэндерами пользователя
type Implementation struct {
	user_v1.UnimplementedUserV1Server
	userService service.UserService
}

// NewImplementation новый экземпляр структуры Implementation
func NewImplementation(noteService service.UserService) *Implementation {
	return &Implementation{
		userService: noteService,
	}
}
