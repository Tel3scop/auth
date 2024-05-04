package auth

import (
	"context"

	authAPI "github.com/Tel3scop/auth/pkg/auth_v1"
)

// Login handler для создания нового пользователя
func (i *Implementation) Login(ctx context.Context, request *authAPI.LoginRequest) (*authAPI.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	return &authAPI.LoginResponse{RefreshToken: *refreshToken}, nil
}
