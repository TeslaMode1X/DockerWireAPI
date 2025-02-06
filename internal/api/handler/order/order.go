package order

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	_ "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/order"
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

// GetUsersOrder
//
// @Summary Get user's order
// @Description Retrieves the current user's order based on the user ID in the context
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} order.Model "User's order details"
// @Failure 401 {object} response.ResponseError "User not logged in"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/order [get]
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

// GetUserOrderByUserID
//
// @Summary Get order by order ID
// @Description Retrieves a specific order by its unique identifier
// @Tags orders
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {object} order.Model "Order details"
// @Failure 400 {object} response.ResponseError "Invalid order ID"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/order/{orderId} [get]
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

// CreateUserOrder
//
// @Summary Create a new user order
// @Description Creates a new order for the current user based on the user ID in the context
// @Tags orders
// @Accept json
// @Produce json
// @Success 201 {string} string "User Order created"
// @Failure 401 {object} response.ResponseError "User not logged in"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/order [post]
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

// AlterUserOrder
//
// @Summary Update user's order
// @Description Updates an existing order by its ID and the user ID in the context
// @Tags orders
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Success 200 {string} string "Order paid successfully"
// @Failure 400 {object} response.ResponseError "Missing or invalid order ID"
// @Failure 401 {object} response.ResponseError "User not logged in"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/order/{orderId} [put]
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

// AddOrderItemIntoOrder
//
// @Summary Add items to user's order
// @Description Adds one or more items to the current user's order
// @Tags orders
// @Accept json
// @Produce json
// @Param request body orderModels.CreateOrderItemRequest true "Order items to add"
// @Success 201 {string} string "Order Items Added"
// @Failure 400 {object} response.ResponseError "Invalid input"
// @Failure 401 {object} response.ResponseError "User not logged in"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/order/order [post]
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
