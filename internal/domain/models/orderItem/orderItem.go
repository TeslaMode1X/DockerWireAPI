package orderItem

import "github.com/gofrs/uuid"

type OrderItem struct {
	BookID   uuid.UUID `json:"book_id"`
	Quantity int       `json:"quantity"`
}

type CreateOrderItemRequest struct {
	Items []OrderItem `json:"items"`
}
