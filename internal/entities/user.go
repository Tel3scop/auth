package entities

import (
	"time"

	"github.com/Tel3scop/auth/pkg/user_v1"
)

// User Структура пользователя
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      user_v1.Role
}
