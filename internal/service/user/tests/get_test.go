package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Tel3scop/auth/internal/client/db/transaction"
	txMocks "github.com/Tel3scop/auth/internal/client/db/transaction/mocks"
	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/repository/mocks"
	"github.com/Tel3scop/auth/internal/service/user"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		user model.User
	}

	var (
		ctx         = context.Background()
		defaultData = model.User{
			ID:        gofakeit.Int64(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Password:  gofakeit.Password(true, true, true, false, false, 8),
			Role:      userAPI.Role(rand.Intn(len(rolesSlice))), // #nosec G404
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
	)

	tests := []struct {
		name                  string
		args                  args
		want                  *model.User
		err                   error
		userMockRepository    func(mc *minimock.Controller) repository.UserRepository
		historyMockRepository func(mc *minimock.Controller) repository.HistoryChangeRepository
		txMockManager         func(mc *minimock.Controller) transaction.TxManager
	}{
		{
			name: "Успешный запуск",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err: nil,
			want: &model.User{
				ID:        defaultData.ID,
				Name:      defaultData.Name,
				Email:     defaultData.Email,
				Role:      defaultData.Role,
				CreatedAt: defaultData.CreatedAt,
				UpdatedAt: defaultData.UpdatedAt,
			},
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				result := &model.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				mock.GetMock.Expect(ctx, defaultData.ID).Return(result, nil)
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				return mock
			},
			txMockManager: func(mc *minimock.Controller) transaction.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "Пользователь не найден",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err:  fmt.Errorf("пользователь не найден"),
			want: nil,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, defaultData.ID).Return(&model.User{}, fmt.Errorf("пользователь не найден"))
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				return mock
			},
			txMockManager: func(mc *minimock.Controller) transaction.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			t.Parallel()

			userRepoMock := tt.userMockRepository(mc)
			historyRepoMock := tt.historyMockRepository(mc)
			txMock := tt.txMockManager(mc)
			newService := user.NewService(userRepoMock, historyRepoMock, txMock)
			actualResult, err := newService.Get(tt.args.ctx, tt.args.user.ID)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, actualResult)
		})
	}
}
