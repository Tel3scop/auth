package user

import (
	"context"
	"fmt"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository/converter"
	"github.com/Tel3scop/helpers/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *serv) Create(ctx context.Context, dto model.User) (int64, error) {
	logger.Info("Creating user...", zap.String("name", dto.Name), zap.String("email", dto.Email))

	var id int64
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("cannot hash password", zap.String("name", dto.Name), zap.String("password", dto.Password), zap.String("email", dto.Email), zap.Error(err))
		return id, fmt.Errorf("не удалось обработать пароль")
	}

	dto.Password = string(hashedPassword)
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.Create(ctx, dto)
		if errTx != nil {
			return errTx
		}

		createdUser, errTx := s.userRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		newHistory, errTx := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, id, createdUser)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.historyChangeRepository.Create(ctx, *newHistory)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	logger.Info("User created", zap.Int64("id", id), zap.String("name", dto.Name), zap.String("email", dto.Email))
	return id, nil
}
