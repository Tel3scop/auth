package user_service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Tel3scop/auth/internal/entities"
	"github.com/Tel3scop/auth/internal/storages/user_storage"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
)

// Create создание нового пользователя. Если пароль и подтверждение не свпадают - возвращаем 0
func Create(ctx context.Context, request *userAPI.CreateRequest) (int64, error) {
	if request.Password != request.PasswordConfirm {
		log.Printf("passwords not equal")

		return 0, fmt.Errorf("passwords not equal")
	}

	now := time.Now()
	userData := entities.User{
		CreatedAt: now,
		UpdatedAt: now,
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		Role:      request.Role,
	}
	createdUser, err := user_storage.Create(ctx, userData)
	if err != nil {

		return 0, err
	}

	return createdUser.ID, nil
}

// GetByID получение пользователя по ID
func GetByID(ctx context.Context, id int64) (entities.User, error) {
	foundedUser, err := user_storage.GetByID(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return entities.User{}, err
	}

	return foundedUser, nil
}

// UpdateByID обновление пользователя по ID
func UpdateByID(ctx context.Context, request *userAPI.UpdateRequest) error {
	updatingData := entities.UpdatingUserData{
		Name:  request.Name,
		Email: request.Email,
		Role:  request.Role,
	}
	_, err := user_storage.UpdateByID(ctx, request.Id, updatingData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// DeleteByID удаление пользователя по ID
func DeleteByID(ctx context.Context, id int64) error {
	err := user_storage.DeleteByID(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
