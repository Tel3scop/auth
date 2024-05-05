package user

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/repository/converter"
	modelRepo "github.com/Tel3scop/auth/internal/repository/user/model"
	"github.com/Tel3scop/helpers/logger"
	"go.uber.org/zap"
)

const (
	tableName = "users"

	columnID        = "id"
	columnName      = "name"
	columnEmail     = "email"
	columnPassword  = "password"
	columnRole      = "role"
	columnCreatedAt = "created_at"
	columnUpdatedAt = "updated_at"
)

type repo struct {
	db db.Client
}

// NewRepository создание репозитория для пользователей
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

// Create Метод создания нового пользователя.
func (r *repo) Create(ctx context.Context, dto model.User) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(columnName, columnPassword, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt).
		Values(dto.Name, dto.Password, dto.Email, dto.Role, dto.CreatedAt, dto.UpdatedAt).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("cannot create query for", zap.String("name", dto.Name), zap.String("password", dto.Password), zap.String("email", dto.Email), zap.Error(err))
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		logger.Error("cannot get result", zap.String("name", dto.Name), zap.String("password", dto.Password), zap.String("email", dto.Email), zap.Error(err))
		return 0, err
	}

	return id, nil
}

// Get получение пользователя по columnID. При его отсутствии возвращаем пустую структуру.
func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(columnID, columnName, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{columnID: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error("cannot get user", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.Error("cannot get data", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return converter.ToUserModelFromRepo(user), nil
}

// GetUserByUsername получение пользователя по columnName. При его отсутствии возвращаем пустую структуру.
func (r *repo) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	builder := sq.Select(columnID, columnName, columnPassword, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{columnName: username}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetPasswordByUsername",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserModelFromRepo(user), nil
}

// Update обновление пользователя по структуре model.UpdatingUserData. При отсутствии пользователя с columnID возвращаем 0.
func (r *repo) Update(ctx context.Context, id int64, data model.UpdatingUserData) (int64, error) {
	builder := sq.Update(tableName).
		Set(columnName, data.Name).
		Set(columnEmail, data.Email).
		Set(columnRole, data.Role).
		Set(columnUpdatedAt, sq.Expr("now()")).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{columnID: id}).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, fmt.Errorf("failed to buildquery : %s", err)
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}
	var userID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		log.Printf("failed to update users: %v", err)
		return 0, fmt.Errorf("failed to update users: %s", err)
	}

	return userID, nil
}

// Delete удаление пользователя по структуре userAPI.UpdateRequest.
func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).PlaceholderFormat(sq.Dollar).Where(sq.Eq{columnID: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("error building delete query: %s", err)
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}
	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("error executing delete query: %s", err)
	}

	if result.RowsAffected() > 0 {
		return nil
	}

	return fmt.Errorf("can not delete record")
}
