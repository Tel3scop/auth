package storages

import (
	"github.com/Tel3scop/auth/internal/storages/users"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Storage контейнер с TableStorages
type Storage struct {
	Users *users.TableStorage
}

// New создать новый контейнер
func New(conn *pgxpool.Pool) *Storage {
	return &Storage{
		Users: users.New(conn),
	}
}
