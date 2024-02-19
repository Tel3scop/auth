package users

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Tel3scop/auth/internal/entities"
	"github.com/Tel3scop/auth/pkg/user_v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	Users     = "users"
	ID        = "id"
	Name      = "name"
	Email     = "email"
	Password  = "password"
	Role      = "role"
	CreatedAt = "created_at"
	UpdatedAt = "updated_at"
)

// TableStorage структура подключения к пулу
type TableStorage struct {
	conn *pgxpool.Pool
}

// New создать новый TableStorage
func New(conn *pgxpool.Pool) *TableStorage {
	return &TableStorage{conn: conn}
}

// Create Метод создания нового пользователя.
func (t TableStorage) Create(ctx context.Context, userData entities.User) (int64, error) {
	builderInsert := sq.Insert(Users).
		PlaceholderFormat(sq.Dollar).
		Columns(Name, Password, Email, Role).
		Values(userData.Name, userData.Password, userData.Email, userData.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return 0, fmt.Errorf("failed to build query: %s", err)
	}

	var userID int64
	err = t.conn.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return 0, fmt.Errorf("failed to insert user: %s", err)
	}

	log.Printf("inserted user with id: %d", userID)
	return userID, nil
}

// GetByID получение пользователя по ID. При его отсутствии возвращаем пустую структуру.
func (t TableStorage) GetByID(ctx context.Context, userID int64) (entities.User, error) {
	builderSelectOne := sq.Select(ID, Name, Email, Role, CreatedAt, UpdatedAt).
		From(Users).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ID: userID}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return entities.User{}, fmt.Errorf("failed to build query: %s", err)
	}

	var (
		id                   int64
		role                 user_v1.Role
		name, email          string
		createdAt, updatedAt time.Time
	)
	err = t.conn.QueryRow(ctx, query, args...).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("failed to select users: %v", err)
		return entities.User{}, fmt.Errorf("failed to select users: %s", err)
	}

	user := entities.User{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Name:      name,
		Email:     email,
		Role:      role,
	}

	return user, nil
}

// UpdateByID обновление пользователя по структуре entities.UpdatingUserData. При отсутствии пользователя с ID возвращаем пустую структуру.
func (t TableStorage) UpdateByID(ctx context.Context, id int64, data entities.UpdatingUserData) (int64, error) {
	builderUpdateOne := sq.Update(Users).
		Set(Name, data.Name).
		Set(Email, data.Email).
		Set(Role, data.Role).
		Set(UpdatedAt, sq.Expr("now()")).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{ID: id}).
		Suffix("RETURNING id")

	query, args, err := builderUpdateOne.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, fmt.Errorf("failed to buildquery : %s", err)
	}

	var userID int64
	err = t.conn.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Printf("failed to update users: %v", err)
		return 0, fmt.Errorf("failed to update users: %s", err)
	}

	return userID, nil
}

// DeleteByID удаление пользователя по структуре userAPI.UpdateRequest.
func (t TableStorage) DeleteByID(ctx context.Context, id int64) error {
	deleteBuilder := sq.Delete(Users).PlaceholderFormat(sq.Dollar).Where(sq.Eq{ID: id})

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("error building delete query: %s", err)
	}

	result, err := t.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing delete query: %s", err)
	}

	if result.RowsAffected() > 0 {
		return nil
	}

	return fmt.Errorf("can not delete record")
}
