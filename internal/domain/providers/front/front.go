package front

import (
	frontHdl "github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/front"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	frontSvc "github.com/TeslaMode1X/DockerWireAPI/internal/service/front"
	"github.com/google/wire"
	"html/template"
	"log/slog"
	"sync"
)

var (
	hdl     *frontHdl.Handler
	hdlOnce sync.Once

	svc     *frontSvc.Service
	svcOnce sync.Once
)

var ProviderSet = wire.NewSet(
	ProvideSetHandler,
	ProvideSetService,
	ProvideSetTemplates,

	wire.Bind(new(interfaces.FrontHandler), new(*frontHdl.Handler)),
	wire.Bind(new(interfaces.FrontService), new(*frontSvc.Service)),
)

func ProvideSetHandler(svc interfaces.FrontService, log *slog.Logger) *frontHdl.Handler {
	hdlOnce.Do(func() {
		hdl = &frontHdl.Handler{
			Svc: svc,
			Log: log,
		}
	})

	return hdl
}

func ProvideSetService(repo interfaces.UserRepository, templates map[string]*template.Template) *frontSvc.Service {
	svcOnce.Do(func() {
		svc = &frontSvc.Service{
			UserRepo:  repo,
			Templates: templates,
		}
	})

	return svc
}

func ProvideSetTemplates() map[string]*template.Template {
	return map[string]*template.Template{
		"main":  template.Must(template.ParseFiles("templates/main.html")),
		"about": template.Must(template.ParseFiles("templates/about.html")),
	}
}
