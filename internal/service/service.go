package service

import (
	"context"

	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/model"
)

// UserService интерфейс для использования в сервисе
type UserService interface {
	Create(ctx context.Context, dto model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, userID int64, data model.UpdatingUserData) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// TxManager интерфейс для создания моков
type TxManager interface {
	// ReadCommitted запускает транзакцию с уровнем изоляции ReadCommitted.
	ReadCommitted(ctx context.Context, f db.Handler) error
}
