package access

import (
	"context"
	"fmt"

	"github.com/Tel3scop/auth/internal/utils"
)

func (s *serv) Check(ctx context.Context, endpoint string, token string) error {

	claims, err := utils.VerifyToken(token, []byte(s.cfg.Encrypt.AccessTokenSecretKey))
	if err != nil {
		return fmt.Errorf("access token is invalid")
	}

	err = s.accessRepository.Check(ctx, endpoint, claims.Role)
	if err != nil {
		return fmt.Errorf("доступ запрещен")
	}

	return nil
}
