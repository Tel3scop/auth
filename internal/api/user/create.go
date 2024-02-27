package user

import (
	"context"
	"fmt"
	"log"

	"github.com/Tel3scop/auth/internal/api/converter"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
)

// Create handler для создания нового пользователя
func (i *Implementation) Create(ctx context.Context, request *userAPI.CreateRequest) (*userAPI.CreateResponse, error) {
	if request.Password != request.PasswordConfirm {
		return nil, fmt.Errorf("passwords not equal")
	}
	id, err := i.userService.Create(ctx, converter.ToUserModelFromRequest(request))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &userAPI.CreateResponse{
		Id: id,
	}, nil
}
