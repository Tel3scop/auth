package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Tel3scop/auth/internal/api/user"
	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/internal/service/mocks"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

var rolesSlice = []int{1, 2}

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx             context.Context
		id              int64
		name            string
		email           string
		password        string
		passwordConfirm string
		role            userAPI.Role
	}

	var (
		ctx         = context.Background()
		defaultData = struct {
			ID       int64
			Name     string
			Email    string
			Password string
			Role     userAPI.Role
		}{
			ID:       gofakeit.Int64(),
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, false, false, 8),
			Role:     userAPI.Role(rand.Intn(len(rolesSlice))), // #nosec G404
		}
	)

	tests := []struct {
		name          string
		args          args
		want          *userAPI.CreateResponse
		err           error
		configureMock func(mc *minimock.Controller, args args) service.UserService
	}{
		{
			name: "Успешный запуск",
			args: args{
				ctx:             ctx,
				id:              defaultData.ID,
				name:            defaultData.Name,
				email:           defaultData.Email,
				password:        defaultData.Password,
				passwordConfirm: defaultData.Password,
				role:            defaultData.Role,
			},
			err:  nil,
			want: &userAPI.CreateResponse{Id: defaultData.ID},
			configureMock: func(mc *minimock.Controller, args args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				dto := model.User{
					Name:     args.name,
					Email:    args.email,
					Password: args.password,
					Role:     args.role,
				}
				mock.CreateMock.Expect(args.ctx, dto).Return(defaultData.ID, nil)
				return mock
			},
		},
		{
			name: "Пароли не совпадают",
			args: args{
				ctx:             context.Background(),
				id:              defaultData.ID,
				name:            defaultData.Name,
				email:           defaultData.Email,
				password:        defaultData.Password,
				passwordConfirm: "another password",
				role:            defaultData.Role,
			},
			want: nil,
			err:  fmt.Errorf("passwords not equal"),
			configureMock: func(mc *minimock.Controller, _ args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				return mock
			},
		},
		{
			name: "Ошибка из сервиса",
			args: args{
				ctx:             context.Background(),
				id:              defaultData.ID,
				name:            defaultData.Name,
				email:           defaultData.Email,
				password:        defaultData.Password,
				passwordConfirm: defaultData.Password,
				role:            defaultData.Role,
			},
			want: nil,
			err:  fmt.Errorf("какая-то ошибка"),
			configureMock: func(mc *minimock.Controller, args args) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				dto := model.User{
					Name:     args.name,
					Email:    args.email,
					Password: args.password,
					Role:     args.role,
				}
				mock.CreateMock.Expect(args.ctx, dto).Return(0, fmt.Errorf("какая-то ошибка"))
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

			request := userAPI.CreateRequest{
				Name:            tt.args.name,
				Email:           tt.args.email,
				Password:        tt.args.password,
				PasswordConfirm: tt.args.passwordConfirm,
				Role:            tt.args.role,
			}
			actualResult, err := api.Create(tt.args.ctx, &request)
			require.Equal(t, tt.err, err)

			require.Equal(t, tt.want, actualResult)
		})
	}
}
