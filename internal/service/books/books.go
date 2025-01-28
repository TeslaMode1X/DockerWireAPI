package books

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type Service struct {
	BookRepo interfaces.BookRepository
}

func (s *Service) GetAllBooks(ctx context.Context) (*[]model.Book, error) {
	const op = "service.books.GetAllBooks"

	books, err := s.BookRepo.GetAllBooks(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return books, nil
}

func (s *Service) GetBookById(ctx context.Context, bookId uuid.UUID) (*model.Book, error) {
	const op = "service.books.GetBookById"

	exists, err := s.BookRepo.IfBookExists(ctx, bookId)
	if !exists {
		return nil, errors.Wrap(errors.New("Book doesn't exists"), op)
	}
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	book, err := s.BookRepo.GetBookById(ctx, bookId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return book, nil
}

func (s *Service) CreateBook(ctx context.Context, book model.Book) (uuid.UUID, error) {
	const op = "service.books.CreateBook"

	createdBookId, err := s.BookRepo.CreateBook(ctx, book)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	return createdBookId, nil
}

func (s *Service) UpdateBookById(ctx context.Context, book model.Book, bookId uuid.UUID) (uuid.UUID, error) {
	const op = "service.books.UpdateBook"

	updatedBookId, err := s.BookRepo.UpdateBookById(ctx, book, bookId)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	return updatedBookId, nil
}

func (s *Service) DeleteBookById(ctx context.Context, bookId uuid.UUID) error {
	const op = "service.books.DeleteBookById"

	err := s.BookRepo.DeleteBookById(ctx, bookId)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
