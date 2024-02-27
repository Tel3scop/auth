package user

import (
	"context"

	"github.com/Tel3scop/auth/internal/api/converter"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update handler для обновления пользователя по ID
func (i *Implementation) Update(ctx context.Context, request *userAPI.UpdateRequest) (*emptypb.Empty, error) {
	_, err := i.userService.Update(ctx, request.GetId(), converter.ToUserModelFromRequestUpdate(request))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
