package books

import (
	"context"
	"database/sql"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) GetAllBooks(ctx context.Context) (*[]model.Book, error) {
	const op = "repository.books.GetAllBooks"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, title, author, price, stock FROM books")
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer stmt.Close()

	var books []model.Book

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	for rows.Next() {
		var book model.Book

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Stock)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}

		books = append(books, book)
	}

	return &books, err
}

func (r *Repository) GetBookById(ctx context.Context, bookId uuid.UUID) (*model.Book, error) {
	const op = "repository.books.GetBookById"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, title, author, price, stock FROM books WHERE id = $1")
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer stmt.Close()

	var book model.Book

	row := stmt.QueryRowContext(ctx, bookId)
	if row.Err() != nil {
		return nil, errors.Wrap(err, op)
	}

	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Stock)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &book, nil
}

func (r *Repository) CreateBook(ctx context.Context, book model.Book) (uuid.UUID, error) {
	const op = "repository.books.CreateBook"

	var bookID uuid.UUID
	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO books (title, author, price, stock, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		book.Title, book.Author, book.Price, book.Stock, time.Now(),
	).Scan(&bookID)

	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	return bookID, nil
}

func (r *Repository) DeleteBookById(ctx context.Context, id uuid.UUID) error {
	const op = "repository.books.DeleteBook"

	stmt, err := r.DB.PrepareContext(ctx, "DELETE FROM books where id = $1")
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repository) UpdateBookById(ctx context.Context, book model.Book, bookId uuid.UUID) (uuid.UUID, error) {
	const op = "repository.books.UpdateBook"

	stmt, err := r.DB.PrepareContext(ctx, `
		UPDATE books 
		SET 
			title = $2, 
			author = $3,
			price = $4,
			stock = $5
		WHERE id = $1
	`)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bookId, book.Title, book.Author, book.Price, book.Stock)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op)
	}

	return bookId, nil
}

func (r *Repository) IfBookExists(ctx context.Context, bookId uuid.UUID) (bool, error) {
	const op = "repository.books.IfBookExists"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT 1 FROM books WHERE id = $1 LIMIT 1")
	if err != nil {
		return false, errors.Wrap(err, op)
	}
	defer stmt.Close()

	var exists bool
	row := stmt.QueryRowContext(ctx, bookId)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(row.Err(), op)
	}

	err = row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, op)
	}

	return true, nil
}
