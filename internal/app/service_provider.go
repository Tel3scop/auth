package app

import (
	"context"
	"log"

	"github.com/Tel3scop/auth/internal/api/user"
	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/client/db/pg"
	"github.com/Tel3scop/auth/internal/client/db/transaction"
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/repository"
	historyChangeRepo "github.com/Tel3scop/auth/internal/repository/history_change"
	userRepo "github.com/Tel3scop/auth/internal/repository/user"
	"github.com/Tel3scop/auth/internal/service"
	userService "github.com/Tel3scop/auth/internal/service/user"
	"github.com/Tel3scop/helpers/closer"
)

type serviceProvider struct {
	config *config.Config

	dbClient                db.Client
	txManager               db.TxManager
	userRepository          repository.UserRepository
	historyChangeRepository repository.HistoryChangeRepository

	userService service.UserService

	userImpl *user.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		cfg, err := config.New()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.Config().Postgres.DSN)
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) HistoryChangeRepository(ctx context.Context) repository.HistoryChangeRepository {
	if s.historyChangeRepository == nil {
		s.historyChangeRepository = historyChangeRepo.NewRepository(s.DBClient(ctx))
	}

	return s.historyChangeRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.HistoryChangeRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}
