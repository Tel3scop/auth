package auth

import (
	"context"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetAccessToken(_ context.Context, refreshToken string) (*string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.cfg.Encrypt.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		Role:     claims.Role,
	},
		[]byte(s.cfg.Encrypt.AccessTokenSecretKey),
		s.cfg.Encrypt.AccessTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}
