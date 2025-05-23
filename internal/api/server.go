package api

import (
	"fmt"
	_ "github.com/TeslaMode1X/DockerWireAPI/docs"
	"github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/auth"
	"github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/books"
	"github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/front"
	"github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/order"
	"github.com/TeslaMode1X/DockerWireAPI/internal/api/handler/user"
	"github.com/TeslaMode1X/DockerWireAPI/internal/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"time"
)

type ServerHTTP struct {
	router http.Handler
}

func NewServeHTTP(cfg *config.Config, authHdl *auth.Handler,
	userHdl *user.Handler, bookHdl *books.Handler,
	frontHdl *front.Handler, orderHdl *order.Handler) *ServerHTTP {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Route("/", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			frontHdl.NewFrontEndHandler(r)
		})
	})

	r.Route("/api/v1/", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			authHdl.NewAuthHandler(r)
			userHdl.NewUserHandler(r)
			bookHdl.NewBookHandler(r)
			orderHdl.NewOrderHandler(r)
		})

		r.Get("/swagger/*", httpSwagger.Handler())
	})

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	handler = cors.AllowAll().Handler(r) // should use if and only if API is open

	return &ServerHTTP{router: handler}
}

func (sh *ServerHTTP) Start(cfg *config.Config, log *slog.Logger) {
	fmt.Print(fmt.Sprintf("Port is %s ", cfg.Server.Port))
	log.Info(fmt.Sprintf("Starting server on port: %s", cfg.Server.Port))
	addr := cfg.Server.Addr + ":" + cfg.Server.Port
	err := http.ListenAndServe(addr, sh.router)
	if err != nil {
		log.Error(err.Error())
		return
	}
}
