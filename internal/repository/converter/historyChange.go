package converter

import (
	"encoding/json"

	"github.com/Tel3scop/auth/internal/model"
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
			return nil, err
		}
	}

	return historyChange, nil
}
