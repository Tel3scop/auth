package user_storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Tel3scop/auth/internal/entities"
	userAPI "github.com/Tel3scop/auth/pkg/user_v1"
)

// SyncMap Эмуляция БД с сиквенсом
type SyncMap struct {
	elems    map[int64]entities.User
	m        sync.RWMutex
	sequence int64
}

var users = &SyncMap{
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

// UpdateByID обновление пользователя по структуре userAPI.UpdateRequest. При его отсутствии возвращаем пустую структуру.
func UpdateByID(ctx context.Context, request *userAPI.UpdateRequest) (entities.User, error) {
	_ = ctx
	var user entities.User
	users.m.Lock()
	defer users.m.Unlock()

	user, ok := users.elems[request.Id]
	if !ok {
		return entities.User{}, fmt.Errorf("user %d not found", request.Id)
	}
	user.Name = request.Name
	user.Email = request.Email
	user.Role = request.Role
	users.elems[request.Id] = user

	return users.elems[request.Id], nil
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
