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

// GetAllBooks
//
// @Summary Get all books
// @Description Retrieves a list of all books available in the system
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} model.Book "List of books"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/book [get]
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

// GetBookById
//
// @Summary Get book by ID
// @Description Retrieves a single book by its unique identifier
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} model.Book "Book details"
// @Failure 400 {object} response.ResponseError "Invalid UUID format"
// @Failure 404 {object} response.ResponseError "Book not found"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/book/{id} [get]
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

// CreateBook
//
// @Summary Create a new book
// @Description Creates a new book and returns the created book ID
// @Tags books
// @Accept json
// @Produce json
// @Param request body model.Book true "Book data"
// @Success 201 {string} string "Created book ID"
// @Failure 400 {object} response.ResponseError "Invalid input"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/book [post]
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

// UpdateBookById
//
// @Summary Update a book by ID
// @Description Updates an existing book by its unique identifier
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param request body model.Book true "Updated book data"
// @Success 200 {string} string "Updated book ID"
// @Failure 400 {object} response.ResponseError "Invalid UUID format or input"
// @Failure 404 {object} response.ResponseError "Book not found"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/book/{id} [put]
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

// DeleteBookById
//
// @Summary Delete a book by ID
// @Description Deletes an existing book by its unique identifier
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {string} string "Book successfully deleted!"
// @Failure 400 {object} response.ResponseError "Invalid UUID format"
// @Failure 404 {object} response.ResponseError "Book not found"
// @Failure 500 {object} response.ResponseError "Internal server error"
// @Router /api/v1/book/{id} [delete]
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
