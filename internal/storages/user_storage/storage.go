package user_storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Tel3scop/auth/internal/entities"
)

// UserDatabase Эмуляция БД с сиквенсом
type UserDatabase struct {
	elems    map[int64]entities.User
	m        sync.RWMutex
	sequence int64
}

var users = &UserDatabase{
	elems: make(map[int64]entities.User),
}

// Create Метод создания нового пользователя.
func Create(ctx context.Context, userData entities.User) (entities.User, error) {
	_ = ctx
	users.m.Lock()
	defer users.m.Unlock()
	users.sequence++
	userData.ID = users.sequence
	users.elems[users.sequence] = userData

	return users.elems[users.sequence], nil
}

// GetByID получение пользователя по ID. При его отсутствии возвращаем пустую структуру.
func GetByID(ctx context.Context, id int64) (entities.User, error) {
	_ = ctx
	users.m.RLock()
	defer users.m.RUnlock()
	user, ok := users.elems[id]
	if !ok {

		return entities.User{}, fmt.Errorf("user %d not found", id)
	}

	return user, nil
}

// UpdateByID обновление пользователя по структуре entities.UpdatingUserData. При отсутствии пользователя с ID возвращаем пустую структуру.
func UpdateByID(ctx context.Context, userID int64, data entities.UpdatingUserData) (entities.User, error) {
	_ = ctx
	var user entities.User
	users.m.Lock()
	defer users.m.Unlock()

	user, ok := users.elems[userID]
	if !ok {

		return entities.User{}, fmt.Errorf("user %d not found", userID)
	}
	user.Name = data.Name
	user.Email = data.Email
	user.Role = data.Role
	users.elems[userID] = user

	return users.elems[userID], nil
}

// DeleteByID удаление пользователя по структуре userAPI.UpdateRequest.
func DeleteByID(ctx context.Context, id int64) error {
	_ = ctx
	users.m.Lock()
	defer users.m.Unlock()
	_, ok := users.elems[id]
	if !ok {

		return fmt.Errorf("user %d not found", id)
	}

	delete(users.elems, id)

	return nil
}
