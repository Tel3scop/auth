package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Tel3scop/auth/internal/api/converter"
	"github.com/Tel3scop/auth/internal/api/user"
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/internal/service/mocks"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx   context.Context
		id    int64
		name  string
		email string
		role  userAPI.Role
	}

	var (
		ctx         = context.Background()
		defaultData = struct {
			ID    int64
			Name  string
			Email string
			Role  userAPI.Role
		}{
			ID:    gofakeit.Int64(),
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
			Role:  userAPI.Role(rand.Intn(len(rolesSlice))), // #nosec G404
		}
	)

	tests := []struct {
		name          string
		args          args
		want          *emptypb.Empty
		err           error
		configureMock func(mc *minimock.Controller, args args) service.UserService
	}{
		{
			name: "Успешный запуск",
			args: args{
				ctx:   ctx,
				id:    defaultData.ID,
				name:  defaultData.Name,
				email: defaultData.Email,
				role:  defaultData.Role,
			},
			err:  nil,
			want: &emptypb.Empty{},
			configureMock: func(mc *minimock.Controller, args args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				request := converter.ToUserModelFromRequestUpdate(&userAPI.UpdateRequest{
					Name:  args.name,
					Email: args.email,
					Role:  args.role,
				})
				mock.UpdateMock.Expect(args.ctx, args.id, request).Return(defaultData.ID, nil)
				return mock
			},
		},
		{
			name: "ID не существует",
			args: args{
				ctx:   context.Background(),
				id:    defaultData.ID,
				name:  defaultData.Name,
				email: defaultData.Email,
				role:  defaultData.Role,
			},
			want: nil,
			err:  fmt.Errorf("not found"),
			configureMock: func(mc *minimock.Controller, args args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				request := converter.ToUserModelFromRequestUpdate(&userAPI.UpdateRequest{
					Id:    args.id,
					Name:  args.name,
					Email: args.email,
					Role:  args.role,
				})
				mock.UpdateMock.Expect(args.ctx, args.id, request).Return(0, fmt.Errorf("not found"))
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			t.Parallel()

			userService := tt.configureMock(mc, tt.args)
			api := user.NewImplementation(userService)

			request := userAPI.UpdateRequest{
				Id:    tt.args.id,
				Name:  tt.args.name,
				Email: tt.args.email,
				Role:  tt.args.role,
			}
			actualResult, err := api.Update(tt.args.ctx, &request)
			require.Equal(t, tt.err, err)

			require.Equal(t, tt.want, actualResult)
		})
	}
}
