package user

import (
	"context"

	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete handler для удаленияпользователя по ID
func (i *Implementation) Delete(ctx context.Context, request *userAPI.DeleteRequest) (*emptypb.Empty, error) {
	err := i.userService.Delete(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
