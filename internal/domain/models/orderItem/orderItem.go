package orderItem

import (
	"github.com/gofrs/uuid"
	"time"
)

type OrderItem struct {
	BookID   uuid.UUID `json:"book_id"`
	Quantity int       `json:"quantity"`
} // @name OrderItemModel

type HistoryOrderItem struct {
	ID         uuid.UUID       `json:"id"`
	TotalPrice float64         `json:"total_price"`
	Status     string          `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	Items      []OrderItemFull `json:"items"`
} // @name HistoryOrderItemModel

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
