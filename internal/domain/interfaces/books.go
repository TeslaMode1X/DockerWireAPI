package interfaces

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/gofrs/uuid"
	"net/http"
)

//go:generate mockery --name BookRepository
type (
	BookRepository interface {
		GetAllBooks(ctx context.Context) (*[]books.Book, error)
		GetBookById(ctx context.Context, bookId uuid.UUID) (*books.Book, error)
		CreateBook(ctx context.Context, book books.Book) (uuid.UUID, error)
		DeleteBookById(ctx context.Context, id uuid.UUID) error
		IfBookExists(ctx context.Context, bookId uuid.UUID) (bool, error)
		UpdateBookById(ctx context.Context, book books.Book, bookId uuid.UUID) (uuid.UUID, error)
	}
)

//go:generate mockery --name BookService
type (
	BookService interface {
		GetAllBooks(ctx context.Context) (*[]books.Book, error)
		GetBookById(ctx context.Context, bookId uuid.UUID) (*books.Book, error)
		CreateBook(ctx context.Context, book books.Book) (uuid.UUID, error)
		DeleteBookById(ctx context.Context, id uuid.UUID) error
		UpdateBookById(ctx context.Context, book books.Book, bookId uuid.UUID) (uuid.UUID, error)
	}
)

//go:generate mockery --name BookHandler
type (
	BookHandler interface {
		GetAllBooks(w http.ResponseWriter, r *http.Request)
		GetBookById(w http.ResponseWriter, r *http.Request)
		CreateBook(w http.ResponseWriter, r *http.Request)
		DeleteBookById(w http.ResponseWriter, r *http.Request)
		UpdateBookById(w http.ResponseWriter, r *http.Request)
	}
)
