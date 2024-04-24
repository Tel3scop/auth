package tests

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Tel3scop/auth/internal/client/db"
	txMocks "github.com/Tel3scop/auth/internal/client/db/transaction/mocks"
	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/repository/converter"
	"github.com/Tel3scop/auth/internal/repository/mocks"
	modelRepo "github.com/Tel3scop/auth/internal/repository/user/model"
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/internal/service/user"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	var rolesSlice = []int{1, 2}

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
		txMockManager         func(mc *minimock.Controller) service.TxManager
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
				request := model.User{
					Name:     defaultData.Name,
					Email:    defaultData.Email,
					Password: defaultData.Password,
					Role:     defaultData.Role,
				}
				mock.CreateMock.Expect(ctx, request).Return(defaultData.ID, nil)
				gettingNewUser := &model.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				mock.GetMock.Expect(ctx, defaultData.ID).Return(gettingNewUser, nil)
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				createdUser := &modelRepo.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				newHistory, _ := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, defaultData.ID, createdUser)
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
			name: "Ошибка создания новой записи",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err:  fmt.Errorf("не удалось создать пользователя"),
			want: 0,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				request := model.User{
					Name:     defaultData.Name,
					Email:    defaultData.Email,
					Password: defaultData.Password,
					Role:     defaultData.Role,
				}
				mock.CreateMock.Expect(ctx, request).Return(0, fmt.Errorf("не удалось создать пользователя"))
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
			name: "Не удалось получить созданного пользователя",
			args: args{
				ctx:  ctx,
				user: defaultData,
			},
			err:  fmt.Errorf("не удалось получить пользователя"),
			want: 0,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				request := model.User{
					Name:     defaultData.Name,
					Email:    defaultData.Email,
					Password: defaultData.Password,
					Role:     defaultData.Role,
				}
				mock.CreateMock.Expect(ctx, request).Return(defaultData.ID, nil)

				mock.GetMock.Expect(ctx, defaultData.ID).Return(&model.User{}, fmt.Errorf("не удалось получить пользователя"))
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
			err:  fmt.Errorf("не удалось создать запись в истории изменений"),
			want: 0,
			userMockRepository: func(mc *minimock.Controller) repository.UserRepository {
				mock := mocks.NewUserRepositoryMock(mc)
				request := model.User{
					Name:     defaultData.Name,
					Email:    defaultData.Email,
					Password: defaultData.Password,
					Role:     defaultData.Role,
				}
				mock.CreateMock.Expect(ctx, request).Return(defaultData.ID, nil)
				gettingNewUser := &model.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				mock.GetMock.Expect(ctx, defaultData.ID).Return(gettingNewUser, nil)
				return mock
			},
			historyMockRepository: func(mc *minimock.Controller) repository.HistoryChangeRepository {
				mock := mocks.NewHistoryChangeRepositoryMock(mc)
				createdUser := &modelRepo.User{
					ID:        defaultData.ID,
					Name:      defaultData.Name,
					Email:     defaultData.Email,
					Role:      defaultData.Role,
					CreatedAt: defaultData.CreatedAt,
					UpdatedAt: defaultData.UpdatedAt,
				}
				newHistory, _ := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, defaultData.ID, createdUser)
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
			request := model.User{
				Name:     tt.args.user.Name,
				Email:    tt.args.user.Email,
				Password: tt.args.user.Password,
				Role:     tt.args.user.Role,
			}
			actualResult, err := newService.Create(tt.args.ctx, request)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, actualResult)
		})
	}
}
