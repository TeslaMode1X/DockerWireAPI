package front

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
)

type Handler struct {
	Svc interfaces.FrontService
	Log *slog.Logger
}

func (h *Handler) NewFrontEndHandler(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.MainPage)
	})
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.MainPage"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	mainPageHTML, err := h.Svc.MainPage(context.Background(), "main")
	if err != nil {
		h.Log.Error("Error in front page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mainPageHTML))
}
