package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Tel3scop/auth/internal/client/db"
	txMocks "github.com/Tel3scop/auth/internal/client/db/transaction/mocks"
	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/repository/converter"
	"github.com/Tel3scop/auth/internal/repository/mocks"
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/internal/service/user"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		user model.User
	}

	var (
		ctx         = context.Background()
		defaultData = model.User{
			ID: gofakeit.Int64(),
		}
	)

	tests := []struct {
		name                  string
		args                  args
		err                   error
		userMockRepository    func(mc *minimock.Controller) repository.UserRepository
		historyMockRepository func(mc *minimock.Controller) repository.HistoryChangeRepository
		txMockManager         func(mc *minimock.Controller) service.TxManager
	}{
		{
			name: "Успешный запуск",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err: nil,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, defaultData.ID).Return(nil)
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				newHistory, _ := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, defaultData.ID, nil)
				mock.CreateMock.Expect(ctx, *newHistory).Return(0, nil)
				return mock
			},
			txMockManager: func(mc *minimock.Controller) service.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "Ошибка удаления записи",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err: fmt.Errorf("не удалось удалить пользователя"),
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, defaultData.ID).Return(fmt.Errorf("не удалось удалить пользователя"))
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				return mock
			},
			txMockManager: func(mc *minimock.Controller) service.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "Не удалось создать запись в истории изменений",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err: fmt.Errorf("не удалось создать запись в истории изменений"),
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, defaultData.ID).Return(nil)
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				newHistory, _ := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, defaultData.ID, nil)
				mock.CreateMock.Expect(ctx, *newHistory).Return(0, fmt.Errorf("не удалось создать запись в истории изменений"))
				return mock
			},
			txMockManager: func(mc *minimock.Controller) service.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
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
			err := newService.Delete(tt.args.ctx, tt.args.user.ID)
			require.Equal(t, tt.err, err)
		})
	}
}
