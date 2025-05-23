package front

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/mainPageParams"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"
	middle "github.com/TeslaMode1X/DockerWireAPI/internal/middleware"
	"github.com/TeslaMode1X/DockerWireAPI/internal/utils/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"strconv"
)

type Handler struct {
	Svc     interfaces.FrontService
	SvcUser interfaces.UserService
	Log     *slog.Logger
}

func (h *Handler) NewFrontEndHandler(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Use(middle.WithOptionalAuth)
		r.Get("/", h.MainPage)

		r.Get("/login", h.LoginPage)
		r.Get("/register", h.RegistrationPage)

		r.Get("/history", h.HistoryPage)

		r.Post("/register/front", h.RegistrationFront)
		r.Post("/login/front", h.LoginFront)

		r.Route("/cart", func(r chi.Router) {
			r.Use(middle.WithAuth)
			r.Get("/add", h.AddCartItems)
			r.Get("/items", h.GetCartItems)
			r.Post("/remove", h.RemoveCartItem)
			r.Get("/success", h.CartCheckout)
		})

		r.Route("/admin", func(r chi.Router) {
			r.Use(middle.WithAuth)
			r.Use(middle.AdminMiddleware)
			r.Get("/", h.AdminPage)

			r.Post("/edit/{id}", h.EditBookFront)
			r.Post("/delete/{id}", h.DeleteBookFront)
		})
	})
}

func (h *Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.MainPage"
	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	var userName string
	if ok {
		user, err := h.SvcUser.GetUserByID(r.Context(), userID)
		if err != nil {
			h.Log.Error("failed to get user by ID", slog.String("error", err.Error()))
		} else {
			userName = user.Username
		}
	}

	roleValue := r.Context().Value("role")
	var role int

	switch v := roleValue.(type) {
	case int:
		role = v
	case float64:
		role = int(v)
	case string:
		var err error
		role, err = strconv.Atoi(v)
		if err != nil {
			h.Log.Error("failed to convert role to int", slog.String("error", err.Error()))
			role = 0
		}
	default:
		h.Log.Error("unexpected type for role", slog.Any("type", fmt.Sprintf("%T", v)))
		role = 0
	}

	fmt.Println(role)

	params := mainPageParams.Model{
		Page:           "main",
		ErrorMessage:   r.URL.Query().Get("error"),
		SuccessMessage: r.URL.Query().Get("success"),
		SearchQuery:    r.URL.Query().Get("search"),
		SortBy:         r.URL.Query().Get("sort"),
		UserName:       userName,
		Role:           role,
	}

	mainPageHTML, err := h.Svc.MainPage(r.Context(), params)
	if err != nil {
		h.Log.Error("Error in front page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mainPageHTML))
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.LoginPage"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	errorMessage := r.URL.Query().Get("error")
	successMessage := r.URL.Query().Get("success")

	loginPage, err := h.Svc.LoginPage(context.Background(), "login", errorMessage, successMessage)
	if err != nil {
		h.Log.Error("Error in front page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(loginPage))
}

func (h *Handler) LoginFront(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.LoginFront"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/login?error=form_error", http.StatusSeeOther)
		return
	}

	err = h.Svc.ProcessLogin(r.Context(), w, r, r.Form)
	if err != nil {
		http.Redirect(w, r, "/login?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/?success=logged_in", http.StatusSeeOther)
}

func (h *Handler) RegistrationPage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.RegistrationPage"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	errorMessage := r.URL.Query().Get("error")
	successMessage := r.URL.Query().Get("success")

	registrationPage, err := h.Svc.RegistrationPage(context.Background(), "registration", errorMessage, successMessage)
	if err != nil {
		h.Log.Error("Error in front page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(registrationPage))
}

func (h *Handler) RegistrationFront(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.RegisterFront"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form error", http.StatusBadRequest)
		return
	}

	err = h.Svc.ProcessRegistration(r.Context(), r.Form)
	if err != nil {
		http.Redirect(w, r, "/register?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login?success=registered", http.StatusSeeOther)
}

func (h *Handler) AdminPage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.AdminPage"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	params := mainPageParams.Model{
		Page:           "admin",
		ErrorMessage:   r.URL.Query().Get("error"),
		SuccessMessage: r.URL.Query().Get("success"),
		SearchQuery:    r.URL.Query().Get("search"),
		SortBy:         r.URL.Query().Get("sort"),
	}

	adminPageHTML, err := h.Svc.AdminPage(r.Context(), params)
	if err != nil {
		h.Log.Error("Error in admin page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(adminPageHTML))
}

func (h *Handler) HistoryPage(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.HistoryPage"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.Log.Error("Error getting userId from context", slog.String("error", ""))
		http.Error(w, "Bad Request Error", http.StatusBadRequest)
		return
	}

	historyPage, err := h.Svc.HistoryPage(r.Context(), userID)
	if err != nil {
		h.Log.Error("Error in history page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(historyPage))
}

func (h *Handler) EditBookFront(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.EditBookFront"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	bookID := chi.URLParam(r, "id")
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/admin?error=form_error", http.StatusSeeOther)
		return
	}

	book := model.Book{
		Title:  r.Form.Get("title"),
		Author: r.Form.Get("author"),
	}

	if price, err := strconv.ParseFloat(r.Form.Get("price"), 64); err == nil {
		book.Price = price
	}
	if stock, err := strconv.Atoi(r.Form.Get("stock")); err == nil {
		book.Stock = stock
	}

	err := h.Svc.EditBook(r.Context(), bookID, &book)
	if err != nil {
		h.Log.Error("failed to edit book", slog.String("error", err.Error()))
		http.Redirect(w, r, "/admin?error=update_failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin?success=book_updated", http.StatusSeeOther)
}

func (h *Handler) DeleteBookFront(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.DeleteBookFront"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	bookID := chi.URLParam(r, "id")

	err := h.Svc.DeleteBook(r.Context(), bookID)
	if err != nil {
		h.Log.Error("failed to delete book", slog.String("error", err.Error()))
		http.Redirect(w, r, "/admin?error=delete_failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin?success=book_deleted", http.StatusSeeOther)
}

// GetCartItems
//
// @Summary Get user's cart items
// @Description Retrieves the list of items in the current user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {array} orderItem.OrderItem "List of cart items"
// @Failure 401 {object} response.ResponseError "User not logged in"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/cart/items [get]
func (h *Handler) GetCartItems(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.GetCartItems"

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

	cartItems, err := h.Svc.GetCartItems(r.Context(), userID)
	if err != nil {
		h.Log.Error("failed to fetch cart items", "error", err)
		http.Error(w, "Failed to load cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cartItems)
}

// AddCartItems
//
// @Summary Add item to cart
// @Description Adds a book to the current user's cart with an optional quantity
// @Tags cart
// @Accept json
// @Produce json
// @Param id query string true "Book ID"
// @Param quantity query integer false "Quantity of the book to add (default: 1)"
// @Success 303 {string} string "Redirects to main page with success message"
// @Failure 400 {object} response.ResponseError "Invalid book ID or missing parameters"
// @Failure 401 {object} response.ResponseError "User not logged in"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/cart/add [get]
func (h *Handler) AddCartItems(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.AddCartItems"

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

	fmt.Println(userID)

	bookIDStr := r.URL.Query().Get("id")
	if bookIDStr == "" {
		h.Log.Error("book ID not provided")
		response.WriteError(w, r, http.StatusBadRequest, errors.New("book ID is required"))
		return
	}

	bookID, err := uuid.FromString(bookIDStr)
	if err != nil {
		h.Log.Error("invalid book ID format", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, errors.New("invalid book ID format"))
		return
	}

	quantityStr := r.URL.Query().Get("quantity")
	quantity := 1

	if quantityStr != "" {
		q, err := strconv.Atoi(quantityStr)
		if err == nil && q > 0 {
			quantity = q
		}
	}

	items := []orderItem.OrderItem{
		{
			BookID:   bookID,
			Quantity: quantity,
		},
	}

	err = h.Svc.AddCartItems(r.Context(), userID, &items)
	if err != nil {
		h.Log.Error("failed to add cart items", "error", err)
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/?success=added_to_cart", http.StatusSeeOther)
}

// RemoveCartItem
//
// @Summary Remove item from cart
// @Description Removes a book from the current user's cart by its ID
// @Tags cart
// @Accept json
// @Produce json
// @Param id query string true "Book ID"
// @Success 303 {string} string "Redirects to main page with success message"
// @Failure 400 {object} response.ResponseError "Invalid book ID or missing parameters"
// @Failure 401 {object} response.ResponseError "User not logged in"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/cart/remove [post]
func (h *Handler) RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.RemoveCartItem"

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

	bookIDStr := r.URL.Query().Get("id")
	if bookIDStr == "" {
		h.Log.Error("book ID not provided")
		response.WriteError(w, r, http.StatusBadRequest, errors.New("book ID is required"))
		return
	}

	bookID, err := uuid.FromString(bookIDStr)
	if err != nil {
		h.Log.Error("invalid book ID format", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, errors.New("invalid book ID format"))
		return
	}

	err = h.Svc.RemoveCartItem(r.Context(), userID, bookID)
	if err != nil {
		h.Log.Error("failed to remove cart item", "error", err)
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/?success=removed_from_cart", http.StatusSeeOther)
}

func (h *Handler) CartCheckout(w http.ResponseWriter, r *http.Request) {
	const op = "handler.front.CartCheckout"
	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.Log.Error("user ID not found in context")
		http.Redirect(w, r, "/login?error=user_not_logged_in", http.StatusSeeOther)
		return
	}

	err := h.Svc.CartCheckout(r.Context(), userID)
	if err != nil {
		h.Log.Error("failed to checkout cart item", "error", err)
		http.Redirect(w, r, "/cart?error=checkout_failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/?success=cart_paid_successfully", http.StatusSeeOther)
}
