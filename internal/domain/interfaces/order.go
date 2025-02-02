package interfaces

import (
	"context"
	orderModel "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/order"
	orderModels "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"
	"github.com/gofrs/uuid"
	"net/http"
)

type (
	OrderRepository interface {
		GetUsersOrder(ctx context.Context, userId string) (*orderModel.Model, error)
		GetUserOrderByUserID(ctx context.Context, orderId string) (*orderModel.Model, error)
		CreateUserOrder(ctx context.Context, userID string) error
		CheckOrderExists(ctx context.Context, userID string) (bool, error)
		AlterUserOrder(ctx context.Context, userID string) error
		AddOrderItemIntoOrder(ctx context.Context, userID string, items *[]orderModels.OrderItem) error
		GetOrderItemsFromOrderID(ctx context.Context, orderID string) (*[]orderModels.OrderItemFull, error)
		RemoveCartItem(ctx context.Context, userID string, bookID uuid.UUID) error
	}

	OrderService interface {
		GetUsersOrder(ctx context.Context, userId string) (*orderModel.Model, error)
		GetUserOrderByUserID(ctx context.Context, orderId string) (*orderModel.Model, error)
		CreateUserOrder(ctx context.Context, userID string) error
		AlterUserOrder(ctx context.Context, userID string) error
		AddOrderItemIntoOrder(ctx context.Context, userID string, bookIDs *[]orderModels.OrderItem) error
	}

	OrderHandler interface {
		GetUsersOrder(w http.ResponseWriter, r *http.Request)
		GetUserOrderByUserID(w http.ResponseWriter, r *http.Request)
		CreateUserOrder(w http.ResponseWriter, r *http.Request)
		AlterUserOrder(w http.ResponseWriter, r *http.Request)
		AddOrderItemIntoOrder(w http.ResponseWriter, r *http.Request)
	}
)
