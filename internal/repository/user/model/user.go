package model

import (
	"time"

	"github.com/Tel3scop/auth/pkg/user_v1"
)

// User Структура пользователя
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     user_v1.Role

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdatingUserData Структура для запроса на изменение пользователя
type UpdatingUserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  user_v1.Role
}
