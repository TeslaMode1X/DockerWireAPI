package order

import (
	"database/sql"
	ordHdl "github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/order"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	ordRepo "github.com/TeslaMode1X/DockerWireAPI/internal/repository/order"
	ordSvc "github.com/TeslaMode1X/DockerWireAPI/internal/service/order"
	"github.com/google/wire"
	"log/slog"
	"sync"
)

var (
	hdl     *ordHdl.Handler
	hdlOnce sync.Once

	svc     *ordSvc.Service
	svcOnce sync.Once

	repo     *ordRepo.Repository
	repoOnce sync.Once
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	ProvideUserHandler,
	ProvideUserService,
	ProvideUserRepository,

	wire.Bind(new(interfaces.OrderHandler), new(*ordHdl.Handler)),
	wire.Bind(new(interfaces.OrderService), new(*ordSvc.Service)),
	wire.Bind(new(interfaces.OrderRepository), new(*ordRepo.Repository)),
)

func ProvideUserHandler(svc interfaces.OrderService, log *slog.Logger) *ordHdl.Handler {
	hdlOnce.Do(func() {
		hdl = &ordHdl.Handler{
			Svc: svc,
			Log: log,
		}
	})

	return hdl
}

func ProvideUserService(repo interfaces.OrderRepository) *ordSvc.Service {
	svcOnce.Do(func() {
		svc = &ordSvc.Service{
			OrderRepo: repo,
		}
	})

	return svc
}

func ProvideUserRepository(db *sql.DB) *ordRepo.Repository {
	repoOnce.Do(func() {
		repo = &ordRepo.Repository{
			DB: db,
		}
	})

	return repo
}
