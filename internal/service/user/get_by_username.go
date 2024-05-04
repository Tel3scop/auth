package user

import (
	"context"

	"github.com/Tel3scop/auth/internal/model"
)

func (s *serv) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
