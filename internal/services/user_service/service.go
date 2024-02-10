package user_service

import (
	"context"
	"log"
	"time"

	"github.com/Tel3scop/auth/internal/entities"
	"github.com/Tel3scop/auth/internal/storages/user_storage"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
)

// Create создание нового пользователя. Если пароль и подтверждение не свпадают - возвращаем 0
func Create(ctx context.Context, request *userAPI.CreateRequest) int64 {
	if request.Password != request.PasswordConfirm {
		log.Printf("passwords not equal")
		return 0
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
		return 0
	}
	return createdUser.ID
}

// GetByID получение пользователя по ID
func GetByID(ctx context.Context, id int64) entities.User {
	foundedUser, err := user_storage.GetByID(ctx, id)
	if err != nil {
		log.Println(err.Error())
	}
	return foundedUser
}

// UpdateByID обновление пользователя по ID
func UpdateByID(ctx context.Context, request *userAPI.UpdateRequest) {
	updatingData := entities.UpdatingUserData{
		Name:  request.Name,
		Email: request.Email,
		Role:  request.Role,
	}
	_, err := user_storage.UpdateByID(ctx, request.Id, updatingData)
	if err != nil {
		log.Println(err.Error())
	}
}

// DeleteByID удаление пользователя по ID
func DeleteByID(ctx context.Context, id int64) {
	err := user_storage.DeleteByID(ctx, id)
	if err != nil {
		log.Println(err.Error())
	}
}
