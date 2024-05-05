package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/helpers/logger"
	"go.uber.org/zap"
)

const (
	tableName = "history_changes"

	columnID        = "id"
	columnEntity    = "entity"
	columnEntityID  = "entity_id"
	columnValue     = "value"
	columnCreatedAt = "created_at"
)

type repo struct {
	db db.Client
}

// NewRepository создание репозитория для истории изменений
func NewRepository(db db.Client) repository.HistoryChangeRepository {
	return &repo{db: db}
}

// Create метод создания новой записи изменений
func (r *repo) Create(ctx context.Context, dto model.HistoryChange) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(columnEntity, columnEntityID, columnValue).
		Values(dto.Entity, dto.EntityID, dto.Value).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("cannot create query", zap.Any("model", dto), zap.Error(err))
		return 0, err
	}

	q := db.Query{
		Name:     "history_change_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		logger.Error("cannot execute query", zap.Any("model", dto), zap.Error(err))
		return 0, err
	}

	return id, nil

}
