package auth

import (
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/service"
)

type serv struct {
	cfg                     *config.Config
	userRepository          repository.UserRepository
	historyChangeRepository repository.HistoryChangeRepository
}

// NewService функция возвращает новый сервис пользователя
func NewService(
	cfg *config.Config,
	userRepository repository.UserRepository,
	historyChangeRepository repository.HistoryChangeRepository,
) service.AuthService {
	return &serv{
		cfg:                     cfg,
		userRepository:          userRepository,
		historyChangeRepository: historyChangeRepository,
	}
}
