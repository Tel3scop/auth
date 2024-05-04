package auth

import (
	"context"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetRefreshToken(_ context.Context, oldRefreshToken string) (*string, error) {
	claims, err := utils.VerifyToken(oldRefreshToken, []byte(s.cfg.Encrypt.RefreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: claims.Username,
		Role:     claims.Role,
	},
		[]byte(s.cfg.Encrypt.RefreshTokenSecretKey),
		s.cfg.Encrypt.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}
