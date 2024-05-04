package auth

import (
	"context"

	authAPI "github.com/Tel3scop/auth/pkg/auth_v1"
)

// GetAccessToken handler для получения access токена
func (i *Implementation) GetAccessToken(ctx context.Context, request *authAPI.GetAccessTokenRequest) (*authAPI.GetAccessTokenResponse, error) {
	accessToken, err := i.authService.GetAccessToken(ctx, request.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &authAPI.GetAccessTokenResponse{AccessToken: *accessToken}, nil
}
