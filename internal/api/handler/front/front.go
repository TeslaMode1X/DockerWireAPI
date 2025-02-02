package front

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/mainPageParams"
	middle "github.com/TeslaMode1X/DockerWireAPI/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"strconv"
)

type Handler struct {
	Svc interfaces.FrontService
	Log *slog.Logger
}

func (h *Handler) NewFrontEndHandler(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.MainPage)

		r.Get("/login", h.LoginPage)
		r.Get("/register", h.RegistrationPage)

		r.Post("/register/front", h.RegistrationFront)
		r.Post("/login/front", h.LoginFront)

		r.Route("/admin", func(r chi.Router) {
			r.Use(middle.WithAuth)
			//r.Use(middle.AdminMiddleware)
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

	params := mainPageParams.Model{
		Page:           "main",
		ErrorMessage:   r.URL.Query().Get("error"),
		SuccessMessage: r.URL.Query().Get("success"),
		SearchQuery:    r.URL.Query().Get("search"),
		SortBy:         r.URL.Query().Get("sort"),
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
		h.Log.Error("failed to edit book", err)
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
		h.Log.Error("failed to delete book", err)
		http.Redirect(w, r, "/admin?error=delete_failed", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/admin?success=book_deleted", http.StatusSeeOther)
}
