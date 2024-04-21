package user

import (
	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/service"
)

type serv struct {
	userRepository          repository.UserRepository
	historyChangeRepository repository.HistoryChangeRepository
	txManager               db.TxManager
}

// NewService функция возвращает новый сервис пользователя
func NewService(
	userRepository repository.UserRepository,
	historyChangeRepository repository.HistoryChangeRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository:          userRepository,
		historyChangeRepository: historyChangeRepository,
		txManager:               txManager,
	}
}
