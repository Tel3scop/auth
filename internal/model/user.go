package model

import (
	"time"

	"github.com/Tel3scop/auth/pkg/user_v1"
)

// Role используется для ограничения доступных ролей
type Role int64

const (
	// RoleUnspecified неиспользуемая роль
	RoleUnspecified Role = iota
	// RoleUser пользователь
	RoleUser
	// RoleAdmin администратор
	RoleAdmin
)

// User Структура пользователя
type User struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	EncryptedPassword []byte `json:"-"`
	Role              Role

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdatingUserData Структура для запрсоа на изменение пользователя
type UpdatingUserData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  user_v1.Role
}

// UserInfo информация о пользователе из токена
type UserInfo struct {
	Username string `json:"username"`
	Role     int64  `json:"role"`
}
