package user

import (
	"context"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository/converter"
)

func (s *serv) Create(ctx context.Context, dto model.User) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
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
