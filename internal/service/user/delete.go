package user

import (
	"context"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository/converter"
)

func (s *serv) Delete(ctx context.Context, requestID int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.userRepository.Delete(ctx, requestID)
		if errTx != nil {
			return errTx
		}

		newHistory, errTx := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, requestID, nil)
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
		return err
	}

	return nil
}
