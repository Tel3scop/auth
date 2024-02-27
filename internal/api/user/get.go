package user

import (
	"context"

	"github.com/Tel3scop/auth/internal/api/converter"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
)

// Get handler для получения нового пользователя по ID
func (i *Implementation) Get(ctx context.Context, request *userAPI.GetRequest) (*userAPI.GetResponse, error) {
	user, err := i.userService.Get(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToUserResponseFromModel(user), nil
}
