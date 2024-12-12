package user

import (
	"database/sql"
	usrHdl "github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/user"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	usrRepo "github.com/TeslaMode1X/DockerWireAPI/internal/repository/user"
	usrSvc "github.com/TeslaMode1X/DockerWireAPI/internal/service/user"
	"github.com/google/wire"
	"log/slog"
	"sync"
)

var (
	hdl     *usrHdl.Handler
	hdlOnce sync.Once

	svc     *usrSvc.Service
	svcOnce sync.Once

	repo     *usrRepo.Repository
	repoOnce sync.Once
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	ProvideUserHandler,
	ProvideUserService,
	ProvideUserRepository,

	wire.Bind(new(interfaces.UserHandler), new(*usrHdl.Handler)),
	wire.Bind(new(interfaces.UserService), new(*usrSvc.Service)),
	wire.Bind(new(interfaces.UserRepository), new(*usrRepo.Repository)),
)

func ProvideUserHandler(svc interfaces.UserService, log *slog.Logger) *usrHdl.Handler {
	hdlOnce.Do(func() {
		hdl = &usrHdl.Handler{
			Svc: svc,
			Log: log,
		}
	})

	return hdl
}

func ProvideUserService(repo interfaces.UserRepository) *usrSvc.Service {
	svcOnce.Do(func() {
		svc = &usrSvc.Service{
			Repo: repo,
		}
	})

	return svc
}

func ProvideUserRepository(db *sql.DB) *usrRepo.Repository {
	repoOnce.Do(func() {
		repo = &usrRepo.Repository{
			DB: db,
		}
	})

	return repo
}
