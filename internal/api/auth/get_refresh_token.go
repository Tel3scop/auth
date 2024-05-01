package auth

import (
	"context"

	authAPI "github.com/Tel3scop/auth/pkg/auth_v1"
)

// GetRefreshToken handler для получения refresh токена
func (i *Implementation) GetRefreshToken(ctx context.Context, request *authAPI.GetRefreshTokenRequest) (*authAPI.GetRefreshTokenResponse, error) {
	token, err := i.authService.GetRefreshToken(ctx, request.GetOldRefreshToken())
	if err != nil {
		return nil, err
	}

	return &authAPI.GetRefreshTokenResponse{RefreshToken: *token}, nil
}
