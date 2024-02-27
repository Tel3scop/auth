package converter

import (
	"time"

	"github.com/Tel3scop/auth/internal/model"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserResponseFromModel функция получения ответа из модели
func ToUserResponseFromModel(user *model.User) *userAPI.GetResponse {
	return &userAPI.GetResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ToUserModelFromRequest функция получения модели пользователя из запроса
func ToUserModelFromRequest(request *userAPI.CreateRequest) model.User {
	return model.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		Role:      request.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ToUserModelFromRequestUpdate функция получения модели пользователя из запроса на изменение
func ToUserModelFromRequestUpdate(request *userAPI.UpdateRequest) model.UpdatingUserData {
	return model.UpdatingUserData{
		Name:  request.Name,
		Email: request.Email,
		Role:  request.Role,
	}
}
