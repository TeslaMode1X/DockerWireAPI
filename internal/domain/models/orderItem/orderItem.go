package orderItem

import "github.com/gofrs/uuid"

type OrderItem struct {
	BookID   uuid.UUID `json:"book_id"`
	Quantity int       `json:"quantity"`
} // @name OrderItemModel

type OrderItemFull struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	BookID   uuid.UUID `json:"book_id"`
	Quantity int       `json:"quantity"`
	Price    float64   `json:"price"`
} // @name OrderItemFullModel

type CreateOrderItemRequest struct {
	Items []OrderItem `json:"items"`
} // @name CreateOrderItemRequestModel
