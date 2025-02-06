package order

import "github.com/gofrs/uuid"

type Model struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"total_price"`
} // @name OrderModel
