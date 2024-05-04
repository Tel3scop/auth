package converter

import (
	"github.com/Tel3scop/auth/internal/model"
	modelRepo "github.com/Tel3scop/auth/internal/repository/user/model"
)

// ToUserModelFromRepo конвертация модели из репозитория в internal/model
func ToUserModelFromRepo(info modelRepo.User) *model.User {
	return &model.User{
		ID:        info.ID,
		Name:      info.Name,
		Email:     info.Email,
		Password:  info.Password,
		Role:      model.Role(info.Role),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}
