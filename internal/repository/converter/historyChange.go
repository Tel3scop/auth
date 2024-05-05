package converter

import (
	"encoding/json"

	"github.com/Tel3scop/auth/internal/model"
	"github.com/Tel3scop/helpers/logger"
	"go.uber.org/zap"
)

// ToHistoryChangeRepoFromEntity конвертер модели в json
func ToHistoryChangeRepoFromEntity(entity string, entityID int64, entityModel interface{}) (*model.HistoryChange, error) {
	historyChange := &model.HistoryChange{
		Entity:   entity,
		EntityID: entityID,
	}

	if entityModel != nil {
		var err error
		historyChange.Value, err = json.Marshal(entityModel)
		if err != nil {
			logger.Error("cannot marshalling", zap.Any("model", entityModel), zap.Error(err))
			return nil, err
		}
	}

	return historyChange, nil
}
