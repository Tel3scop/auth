package user

import (
	"context"
	"fmt"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository/converter"
	"golang.org/x/crypto/bcrypt"
)

func (s *serv) Create(ctx context.Context, dto model.User) (int64, error) {
	var id int64
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("cannot hash password %s:%s %s", dto.Name, dto.Password, err)
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

	return id, nil
}
