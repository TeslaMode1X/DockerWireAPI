package interfaces

import (
	"context"
	orderModel "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/order"
	orderModels "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"
	"github.com/gofrs/uuid"
	"net/http"
)

//go:generate mockery --name OrderRepository
type (
	OrderRepository interface {
		GetOrdersByUserID(ctx context.Context, userID string) ([]orderModels.HistoryOrderItem, error)
		GetUsersOrder(ctx context.Context, userId string) (*orderModel.Model, error)
		GetUserOrderByUserID(ctx context.Context, orderId string) (*orderModel.Model, error)
		ChangeStatusOfCart(ctx context.Context, userId, orderId string) error
		CreateUserOrder(ctx context.Context, userID string) error
		CheckOrderExists(ctx context.Context, userID string) (bool, error)
		AlterUserOrder(ctx context.Context, userID string) error
		AddOrderItemIntoOrder(ctx context.Context, userID string, items *[]orderModels.OrderItem) error
		GetOrderItemsFromOrderID(ctx context.Context, orderID string) (*[]orderModels.OrderItemFull, error)
		RemoveCartItem(ctx context.Context, userID string, bookID uuid.UUID) error
		AlterUserOrderByID(ctx context.Context, userID, orderID uuid.UUID) error
	}
)

//go:generate mockery --name OrderService
type (
	OrderService interface {
		GetUsersOrder(ctx context.Context, userId string) (*orderModel.Model, error)
		GetUserOrderByUserID(ctx context.Context, orderId string) (*orderModel.Model, error)
		CreateUserOrder(ctx context.Context, userID string) error
		AlterUserOrder(ctx context.Context, userID string) error
		AddOrderItemIntoOrder(ctx context.Context, userID string, bookIDs *[]orderModels.OrderItem) error
		AlterUserOrderByID(ctx context.Context, userID, orderID string) error
	}
)

//go:generate mockery --name OrderHandler
type (
	OrderHandler interface {
		GetUsersOrder(w http.ResponseWriter, r *http.Request)
		GetUserOrderByUserID(w http.ResponseWriter, r *http.Request)
		CreateUserOrder(w http.ResponseWriter, r *http.Request)
		AlterUserOrder(w http.ResponseWriter, r *http.Request)
		AddOrderItemIntoOrder(w http.ResponseWriter, r *http.Request)
	}
)
