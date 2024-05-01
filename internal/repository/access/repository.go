package access

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/repository"
)

const (
	tableName = "accesses"

	columnID     = "id"
	columnName   = "name"
	columnRoleID = "role_id"
)

type repo struct {
	db db.Client
}

// NewRepository создание репозитория для пользователей
func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}

// Check Метод создания нового пользователя.
func (r *repo) Check(ctx context.Context, endpoint string, roleID int64) error {
	builder := sq.Select(columnID).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{columnRoleID: roleID}).
		Where(sq.Eq{columnName: endpoint})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_access.Check",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return err
	}

	return nil

}
