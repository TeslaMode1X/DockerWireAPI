package books

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/TeslaMode1X/DockerWireAPI/internal/utils/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"log/slog"
	"net/http"
)

type Handler struct {
	Svc interfaces.BookService
	Log *slog.Logger
}

func (h *Handler) NewBookHandler(r chi.Router) {
	r.Route("/book", func(r chi.Router) {
		r.Get("/", h.GetAllBooks)
		r.Get("/{id}", h.GetBookById)
		r.Post("/", h.CreateBook)
		r.Put("/{id}", h.UpdateBookById)
		r.Delete("/{id}", h.DeleteBookById)
	})
}

func (h *Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	const op = "handler.books.GetAllBooks"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	books, err := h.Svc.GetAllBooks(context.Background())
	if err != nil {
		h.Log.Error("error getting all books", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusOK, books)
}

func (h *Handler) GetBookById(w http.ResponseWriter, r *http.Request) {
	const op = "handler.books.GetBookById"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var id = chi.URLParam(r, "id")

	parsedID, err := uuid.FromString(id)
	if err != nil {
		h.Log.Error("failed to parse UUID", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	book, err := h.Svc.GetBookById(context.Background(), parsedID)
	if err != nil {
		h.Log.Error("error getting book", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusOK, book)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	const op = "handler.books.CreateBook"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var newBook model.Book
	if err := render.DecodeJSON(r.Body, &newBook); err != nil {
		h.Log.Error("failed to decode request body", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.Svc.CreateBook(context.Background(), newBook)
	if err != nil {
		h.Log.Error("error creating book", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusCreated, id)
}

func (h *Handler) UpdateBookById(w http.ResponseWriter, r *http.Request) {
	const op = "handler.books.UpdateBookById"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var id = chi.URLParam(r, "id")

	parsedID, err := uuid.FromString(id)
	if err != nil {
		h.Log.Error("failed to parse UUID", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	var newBook model.Book
	if err := render.DecodeJSON(r.Body, &newBook); err != nil {
		h.Log.Error("failed to decode request body", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	updateId, err := h.Svc.UpdateBookById(context.Background(), newBook, parsedID)
	if err != nil {
		h.Log.Error("error updating book", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusOK, updateId)
}

func (h *Handler) DeleteBookById(w http.ResponseWriter, r *http.Request) {
	const op = "handler.books.DeleteBookById"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var id = chi.URLParam(r, "id")

	parsedID, err := uuid.FromString(id)
	if err != nil {
		h.Log.Error("failed to parse UUID", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, err)
		return
	}

	err = h.Svc.DeleteBookById(context.Background(), parsedID)
	if err != nil {
		h.Log.Error("error deleting book", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusOK, "book successfully deleted!")
}
