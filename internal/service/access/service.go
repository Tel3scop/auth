package access

import (
	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/repository"
	"github.com/Tel3scop/auth/internal/service"
)

type serv struct {
	accessRepository        repository.AccessRepository
	historyChangeRepository repository.HistoryChangeRepository
	txManager               db.TxManager
	cfg                     *config.Config
}

// NewService функция возвращает новый сервис пользователя
func NewService(
	cfg *config.Config,
	accessRepository repository.AccessRepository,
	historyChangeRepository repository.HistoryChangeRepository,
	txManager db.TxManager,
) service.AccessService {
	return &serv{
		cfg:                     cfg,
		accessRepository:        accessRepository,
		historyChangeRepository: historyChangeRepository,
		txManager:               txManager,
	}
}
