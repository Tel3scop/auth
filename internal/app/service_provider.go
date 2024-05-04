package app

import (
	"context"
	"log"

	"github.com/Tel3scop/auth/internal/api/access"
	"github.com/Tel3scop/auth/internal/api/auth"
	"github.com/Tel3scop/auth/internal/api/user"
	"github.com/Tel3scop/auth/internal/client/db"
	"github.com/Tel3scop/auth/internal/client/db/pg"
	"github.com/Tel3scop/auth/internal/client/db/transaction"
	"github.com/Tel3scop/auth/internal/closer"
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/repository"
	accessRepo "github.com/Tel3scop/auth/internal/repository/access"
	historyChangeRepo "github.com/Tel3scop/auth/internal/repository/history_change"
	userRepo "github.com/Tel3scop/auth/internal/repository/user"
	"github.com/Tel3scop/auth/internal/service"
	accessService "github.com/Tel3scop/auth/internal/service/access"
	authService "github.com/Tel3scop/auth/internal/service/auth"
	userService "github.com/Tel3scop/auth/internal/service/user"
)

type serviceProvider struct {
	config *config.Config

	dbClient                db.Client
	txManager               db.TxManager
	accessRepository        repository.AccessRepository
	userRepository          repository.UserRepository
	historyChangeRepository repository.HistoryChangeRepository

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userImpl   *user.Implementation
	accessImpl *access.Implementation
	authImpl   *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		cfg, err := config.New()
		if err != nil {
			log.Fatalf("failed to get config: %s", err.Error())
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

func (s *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if s.accessRepository == nil {
		s.accessRepository = accessRepo.NewRepository(s.DBClient(ctx))
	}

	return s.accessRepository
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

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.Config(),
			s.UserRepository(ctx),
			s.HistoryChangeRepository(ctx),
		)
	}

	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewService(
			s.Config(),
			s.AccessRepository(ctx),
			s.HistoryChangeRepository(ctx),
		)
	}

	return s.accessService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *access.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = access.NewImplementation(s.Config(), s.AccessService(ctx))
	}

	return s.accessImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
