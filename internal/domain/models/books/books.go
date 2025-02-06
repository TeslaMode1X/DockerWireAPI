package books

import "github.com/gofrs/uuid"

type Book struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Author string    `json:"author"`
	Price  float64   `json:"price"`
	Stock  int       `json:"stock"`
} // @name BookModel
