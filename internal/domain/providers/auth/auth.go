package auth

import (
	"database/sql"
	authHdl "github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/auth"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	authRepo "github.com/TeslaMode1X/DockerWireAPI/internal/repository/auth"
	authSvc "github.com/TeslaMode1X/DockerWireAPI/internal/service/auth"
	"github.com/google/wire"
	"log/slog"
	"sync"
)

var (
	hdl     *authHdl.Handler
	hdlOnce sync.Once

	svc     *authSvc.Service
	svcOnce sync.Once

	repo     *authRepo.Repository
	repoOnce sync.Once
)

var ProviderSet = wire.NewSet(
	ProvideSetHandler,
	ProvideSetService,
	ProvideSetRepository,

	wire.Bind(new(interfaces.AuthHandler), new(*authHdl.Handler)),
	wire.Bind(new(interfaces.AuthService), new(*authSvc.Service)),
	wire.Bind(new(interfaces.AuthRepository), new(*authRepo.Repository)),
)

func ProvideSetHandler(svc interfaces.AuthService, log *slog.Logger) *authHdl.Handler {
	hdlOnce.Do(func() {
		hdl = &authHdl.Handler{
			Svc: svc,
			Log: log,
		}
	})

	return hdl
}

func ProvideSetService(repo interfaces.AuthRepository, userRepo interfaces.UserRepository) *authSvc.Service {
	svcOnce.Do(func() {
		svc = &authSvc.Service{
			AuthRepo: repo,
			UserRepo: userRepo,
		}
	})

	return svc
}

func ProvideSetRepository(db *sql.DB) *authRepo.Repository {
	repoOnce.Do(func() {
		repo = &authRepo.Repository{
			DB: db,
		}
	})

	return repo
}
