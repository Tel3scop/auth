package auth

import (
	"context"
	"fmt"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, username, password string) (*string, error) {
	user, errTx := s.userRepository.GetUserByUsername(ctx, username)
	if errTx != nil {
		return nil, fmt.Errorf("пользователь не найден")
	}

	isValidPassword := utils.VerifyPassword(user.Password, password)
	if !isValidPassword {
		return nil, fmt.Errorf("неверный пароль")
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Username: username,
		Role:     int64(user.Role),
	},
		[]byte(s.cfg.Encrypt.RefreshTokenSecretKey),
		s.cfg.Encrypt.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось сгенерировать токен")
	}

	return &refreshToken, nil
}
