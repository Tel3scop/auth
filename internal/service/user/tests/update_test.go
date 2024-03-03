package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/client/db/transaction"
	txMocks "github.com/Tel3scop/auth/internal/client/db/transaction/mocks"
	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/repository/converter"
	"github.com/Tel3scop/auth/internal/repository/mocks"
	modelRepo "github.com/Tel3scop/auth/internal/repository/user/model"
	"github.com/Tel3scop/auth/internal/service/user"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
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
		want                  int64
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
			err:  nil,
			want: defaultData.ID,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				request := model.UpdatingUserData{
					Name:  defaultData.Name,
					Email: defaultData.Email,
					Role:  defaultData.Role,
				}
				mock.UpdateMock.Expect(ctx, defaultData.ID, request).Return(defaultData.ID, nil)
				gettingUser := &model.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				mock.GetMock.Expect(ctx, defaultData.ID).Return(gettingUser, nil)
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				updatedUser := &modelRepo.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				newHistory, _ := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, defaultData.ID, updatedUser)
				mock.CreateMock.Expect(ctx, *newHistory).Return(0, nil)
				return mock
			},
			txMockManager: func(mc *minimock.Controller) transaction.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "Ошибка обновления записи",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err:  fmt.Errorf("не удалось обновить пользователя"),
			want: 0,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				request := model.UpdatingUserData{
					Name:  defaultData.Name,
					Email: defaultData.Email,
					Role:  defaultData.Role,
				}
				mock.UpdateMock.Expect(ctx, defaultData.ID, request).Return(0, fmt.Errorf("не удалось обновить пользователя"))
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				return mock
			},
			txMockManager: func(mc *minimock.Controller) transaction.TxManager {
				mock := txMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "Ошибка получения записи",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err:  fmt.Errorf("не удалось получить обновленного пользователя"),
			want: 0,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				request := model.UpdatingUserData{
					Name:  defaultData.Name,
					Email: defaultData.Email,
					Role:  defaultData.Role,
				}
				mock.UpdateMock.Expect(ctx, defaultData.ID, request).Return(defaultData.ID, nil)
				mock.GetMock.Expect(ctx, defaultData.ID).Return(nil, fmt.Errorf("не удалось получить обновленного пользователя"))
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				return mock
			},
			txMockManager: func(mc *minimock.Controller) transaction.TxManager {
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
			err:  fmt.Errorf("не удалось создать запись в истории изменений"),
			want: 0,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				request := model.UpdatingUserData{
					Name:  defaultData.Name,
					Email: defaultData.Email,
					Role:  defaultData.Role,
				}
				mock.UpdateMock.Expect(ctx, defaultData.ID, request).Return(defaultData.ID, nil)
				gettingUser := &model.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				mock.GetMock.Expect(ctx, defaultData.ID).Return(gettingUser, nil)
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				updatedUser := &modelRepo.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				newHistory, _ := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, defaultData.ID, updatedUser)
				mock.CreateMock.Expect(ctx, *newHistory).Return(0, fmt.Errorf("не удалось создать запись в истории изменений"))
				return mock
			},
			txMockManager: func(mc *minimock.Controller) transaction.TxManager {
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
			request := model.UpdatingUserData{
				Name:  tt.args.user.Name,
				Email: tt.args.user.Email,
				Role:  tt.args.user.Role,
			}
			actualResult, err := newService.Update(tt.args.ctx, tt.args.user.ID, request)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, actualResult)
		})
	}
}
