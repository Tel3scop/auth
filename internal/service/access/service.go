package access

import (
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/service"
)

type serv struct {
	accessRepository        repository.AccessRepository
	historyChangeRepository repository.HistoryChangeRepository
	cfg                     config.Provider
}

// NewService функция возвращает новый сервис пользователя
func NewService(
	cfg config.Provider,
	accessRepository repository.AccessRepository,
	historyChangeRepository repository.HistoryChangeRepository,
) service.AccessService {
	return &serv{
		cfg:                     cfg,
		accessRepository:        accessRepository,
		historyChangeRepository: historyChangeRepository,
	}
}
