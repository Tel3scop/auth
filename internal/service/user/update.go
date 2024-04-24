package user

import (
	"context"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/auth/internal/repository/converter"
)

func (s *serv) Update(ctx context.Context, requestID int64, data model.UpdatingUserData) (int64, error) {
	var userID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		userID, errTx = s.userRepository.Update(ctx, requestID, data)
		if errTx != nil {
			return errTx
		}

		userModel, errTx := s.userRepository.Get(ctx, userID)
		if errTx != nil {
			return errTx
		}

		newHistory, errTx := converter.ToHistoryChangeRepoFromEntity(model.EntityUser, requestID, userModel)
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

	return userID, nil
}
