package repository

import (
	"context"

	"github.com/Tel3scop/auth/internal/model"
)

// UserRepository интерфейс репозитория пользователей
type UserRepository interface {
	Create(ctx context.Context, dto model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, userID int64, data model.UpdatingUserData) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// HistoryChangeRepository интерфейс репозитория истории изменения
type HistoryChangeRepository interface {
	Create(ctx context.Context, dto model.HistoryChange) (int64, error)
}
