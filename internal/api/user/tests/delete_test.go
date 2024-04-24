package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Tel3scop/auth/internal/api/user"
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/internal/service/mocks"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx         = context.Background()
		defaultData = struct{ ID int64 }{ID: gofakeit.Int64()}
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
				ctx: ctx,
				id:  defaultData.ID,
			},
			err:  nil,
			want: &emptypb.Empty{},
			configureMock: func(mc *minimock.Controller, args args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(args.ctx, args.id).Return(nil)
				return mock
			},
		},
		{
			name: "ID не существует",
			args: args{
				ctx: ctx,
				id:  defaultData.ID,
			},
			err: fmt.Errorf("not found"),
			configureMock: func(mc *minimock.Controller, args args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(args.ctx, args.id).Return(fmt.Errorf("not found"))
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

			request := userAPI.DeleteRequest{Id: tt.args.id}
			actualResult, err := api.Delete(tt.args.ctx, &request)
			require.Equal(t, tt.err, err)

			require.Equal(t, tt.want, actualResult)
		})
	}
}
