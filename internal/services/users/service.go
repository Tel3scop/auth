package users

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/entities"
	"github.com/Tel3scop/auth/internal/storages"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
)

type storage interface {
	Create(ctx context.Context, userData entities.User) (int64, error)
	GetByID(ctx context.Context, id int64) (entities.User, error)
	UpdateByID(ctx context.Context, userID int64, data entities.UpdatingUserData) (int64, error)
	DeleteByID(ctx context.Context, id int64) error
}

// Service структура сервиса
type Service struct {
	usersStorage storage
	config       *config.Config
}

// NewService создать новый сервис
func NewService(cfg *config.Config, storages *storages.Storage) Service {
	return Service{
		config:       cfg,
		usersStorage: storages.Users,
	}
}

// Create создание нового пользователя. Если пароль и подтверждение не свпадают - возвращаем 0
func (s *Service) Create(ctx context.Context, request *userAPI.CreateRequest) (int64, error) {
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
	createdUserID, err := s.usersStorage.Create(ctx, userData)
	if err != nil {

		return 0, err
	}

	return createdUserID, nil
}

// GetByID получение пользователя по ID
func (s *Service) GetByID(ctx context.Context, id int64) (entities.User, error) {
	foundedUser, err := s.usersStorage.GetByID(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return entities.User{}, err
	}

	return foundedUser, nil
}

// UpdateByID обновление пользователя по ID
func (s *Service) UpdateByID(ctx context.Context, request *userAPI.UpdateRequest) error {
	updatingData := entities.UpdatingUserData{
		Name:  request.Name,
		Email: request.Email,
		Role:  request.Role,
	}
	_, err := s.usersStorage.UpdateByID(ctx, request.Id, updatingData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// DeleteByID удаление пользователя по ID
func (s *Service) DeleteByID(ctx context.Context, id int64) error {
	err := s.usersStorage.DeleteByID(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
