package model

import "github.com/dgrijalva/jwt-go"

// UserClaims claims пользователя
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     int64  `json:"role"`
}
