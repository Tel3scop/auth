package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Tel3scop/auth/internal/api/converter"
	"github.com/Tel3scop/auth/internal/api/user"
	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/internal/service/mocks"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		id        int64
		name      string
		email     string
		role      userAPI.Role
		createdAt time.Time
		updatedAt time.Time
	}

	var (
		ctx         = context.Background()
		defaultData = struct {
			ID        int64
			Name      string
			Email     string
			Role      userAPI.Role
			CreatedAt time.Time
			UpdatedAt time.Time
		}{
			ID:        gofakeit.Int64(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Role:      userAPI.Role(rand.Intn(len(rolesSlice))), // #nosec G404
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
	)

	tests := []struct {
		name          string
		args          args
		want          *userAPI.GetResponse
		err           error
		configureMock func(mc *minimock.Controller, args args) service.UserService
	}{
		{
			name: "Успешный запуск",
			args: args{
				ctx:       ctx,
				id:        defaultData.ID,
				name:      defaultData.Name,
				email:     defaultData.Email,
				role:      defaultData.Role,
				createdAt: defaultData.CreatedAt,
				updatedAt: defaultData.UpdatedAt,
			},
			err: nil,
			want: converter.ToUserResponseFromModel(&model.User{
				ID:        defaultData.ID,
				Name:      defaultData.Name,
				Email:     defaultData.Email,
				Role:      defaultData.Role,
				CreatedAt: defaultData.CreatedAt,
				UpdatedAt: defaultData.UpdatedAt,
			}),
			configureMock: func(mc *minimock.Controller, args args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(args.ctx, args.id).Return(&model.User{
					ID:        args.id,
					Name:      args.name,
					Email:     args.email,
					Role:      args.role,
					CreatedAt: args.createdAt,
					UpdatedAt: args.updatedAt,
				}, nil)
				return mock
			},
		},
		{
			name: "ID не существует",
			args: args{
				ctx:       ctx,
				id:        defaultData.ID,
				name:      defaultData.Name,
				email:     defaultData.Email,
				role:      defaultData.Role,
				createdAt: defaultData.CreatedAt,
				updatedAt: defaultData.UpdatedAt,
			},
			err: fmt.Errorf("not found"),
			configureMock: func(mc *minimock.Controller, args args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(args.ctx, args.id).Return(nil, fmt.Errorf("not found"))
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

			request := userAPI.GetRequest{Id: tt.args.id}
			actualResult, err := api.Get(tt.args.ctx, &request)
			require.Equal(t, tt.err, err)

			require.Equal(t, tt.want, actualResult)
		})
	}
}
