package access

import (
	"github.com/Tel3scop/auth/internal/config"
	"github.com/Tel3scop/auth/internal/service"
	"github.com/Tel3scop/auth/pkg/access_v1"
)

// Implementation структура для работы с хэндерами пользователя
type Implementation struct {
	access_v1.UnimplementedAccessV1Server
	accessService service.AccessService
	cfg           *config.Config
}

// NewImplementation новый экземпляр структуры Implementation
func NewImplementation(cfg *config.Config, accessService service.AccessService) *Implementation {
	return &Implementation{
		cfg:           cfg,
		accessService: accessService,
	}
}
