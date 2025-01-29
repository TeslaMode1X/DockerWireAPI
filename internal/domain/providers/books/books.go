package books

import (
	"database/sql"
	bookHdl "github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/books"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	bookRepo "github.com/TeslaMode1X/DockerWireAPI/internal/repository/books"
	bookSvc "github.com/TeslaMode1X/DockerWireAPI/internal/service/books"
	"github.com/google/wire"
	"log/slog"
	"sync"
)

var (
	hdl     *bookHdl.Handler
	hdlOnce sync.Once

	svc     *bookSvc.Service
	svcOnce sync.Once

	repo     *bookRepo.Repository
	repoOnce sync.Once
)

var ProviderSet = wire.NewSet(
	ProvideSetHandler,
	ProvideSetService,
	ProvideSetRepository,

	wire.Bind(new(interfaces.BookHandler), new(*bookHdl.Handler)),
	wire.Bind(new(interfaces.BookService), new(*bookSvc.Service)),
	wire.Bind(new(interfaces.BookRepository), new(*bookRepo.Repository)),
)

func ProvideSetHandler(svc interfaces.BookService, log *slog.Logger) *bookHdl.Handler {
	hdlOnce.Do(func() {
		hdl = &bookHdl.Handler{
			Svc: svc,
			Log: log,
		}
	})

	return hdl
}

func ProvideSetService(repo interfaces.BookRepository) *bookSvc.Service {
	svcOnce.Do(func() {
		svc = &bookSvc.Service{
			BookRepo: repo,
		}
	})

	return svc
}

func ProvideSetRepository(db *sql.DB) *bookRepo.Repository {
	repoOnce.Do(func() {
		repo = &bookRepo.Repository{
			DB: db,
		}
	})

	return repo
}
