package model

import (
	"time"
)

const (
	// EntityUser - сущность пользователя
	EntityUser = "user"
)

// HistoryChange Структура записи истории изменений
type HistoryChange struct {
	ID       int64  `json:"id"`
	Entity   string `json:"entity"`
	EntityID int64  `json:"entity_id"`
	Value    []byte `json:"value"`

	CreatedAt time.Time `json:"created_at"`
}
