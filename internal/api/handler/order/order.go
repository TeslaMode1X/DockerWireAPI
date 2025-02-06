package order

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	orderModels "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"
	"github.com/TeslaMode1X/DockerWireAPI/internal/utils/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
)

type Handler struct {
	Log *slog.Logger
	Svc interfaces.OrderService
}

func (h *Handler) NewOrderHandler(r chi.Router) {
	r.Route("/order", func(r chi.Router) {
		r.Get("/", h.GetUsersOrder)

		r.Get("/{orderId}", h.GetUserOrderByUserID)

		r.Post("/", h.CreateUserOrder)

		r.Post("/order", h.AddOrderItemIntoOrder)

		r.Put("/{orderId}", h.AlterUserOrder)
	})
}

func (h *Handler) GetUsersOrder(w http.ResponseWriter, r *http.Request) {
	const op = "handler.order.GetUsersOrder"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.Log.Error("user ID not found in context")
		response.WriteError(w, r, http.StatusUnauthorized, errors.New("user not logged in"))
		return
	}

	order, err := h.Svc.GetUsersOrder(context.Background(), userID)
	if err != nil {
		h.Log.Error("failed to get users order", "error", err)
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusOK, order)
}

func (h *Handler) GetUserOrderByUserID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.order.GetUserOrderByUserID"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	orderID := chi.URLParam(r, "orderId")

	order, err := h.Svc.GetUserOrderByUserID(context.Background(), orderID)
	if err != nil {
		h.Log.Error("failed to get users order", "error", err)
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusOK, order)
}

func (h *Handler) CreateUserOrder(w http.ResponseWriter, r *http.Request) {
	const op = "handler.order.CreateUserOrder"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.Log.Error("user ID not found in context")
		response.WriteError(w, r, http.StatusUnauthorized, errors.New("user not logged in"))
		return
	}

	err := h.Svc.CreateUserOrder(context.Background(), userID)
	if err != nil {
		h.Log.Error("failed to create order", "error", err)
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusCreated, "User Order created")
}

func (h *Handler) AlterUserOrder(w http.ResponseWriter, r *http.Request) {
	const op = "handler.order.AlterUserOrder"
	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.Log.Error("user ID not found in context")
		response.WriteError(w, r, http.StatusUnauthorized, errors.New("user not logged in"))
		return
	}

	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		h.Log.Error("order ID not found in URL param")
		response.WriteError(w, r, http.StatusBadRequest, errors.New("missing orderId"))
		return
	}

	err := h.Svc.AlterUserOrderByID(r.Context(), userID, orderID)
	if err != nil {
		h.Log.Error("failed to alter order", "error", err)
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusOK, "Order paid successfully")
}

func (h *Handler) AddOrderItemIntoOrder(w http.ResponseWriter, r *http.Request) {
	const op = "handler.order.AddOrderItemIntoOrder"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.Log.Error("user ID not found in context")
		response.WriteError(w, r, http.StatusUnauthorized, errors.New("user not logged in"))
		return
	}

	var req orderModels.CreateOrderItemRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.Log.Error("failed to decode request body", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	err := h.Svc.AddOrderItemIntoOrder(r.Context(), userID, &req.Items)
	if err != nil {
		h.Log.Error("failed to add order item to order", "error", err)
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusCreated, "Order Items Added")
}
