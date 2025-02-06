package user

import (
	"context"
	"errors"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	"github.com/TeslaMode1X/DockerWireAPI/internal/service"
	"github.com/TeslaMode1X/DockerWireAPI/internal/utils/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
	"log/slog"
	"net/http"
)

type Handler struct {
	Svc interfaces.UserService
	Log *slog.Logger
}

func (h *Handler) NewUserHandler(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/{id}", h.GetUserByID)
	})
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.user.GetUserByID"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var id = chi.URLParam(r, "id")
	_, err := uuid.FromString(id)
	if err != nil {
		h.Log.Error("failed to parse UUID", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	result, err := h.Svc.GetUserByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.WriteError(w, r, http.StatusNotFound, err)
			return
		}
		response.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	response.WriteJson(w, r, http.StatusOK, result)
}
