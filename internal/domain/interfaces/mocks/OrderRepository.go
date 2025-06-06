// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	order "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/order"

	orderItem "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"

	uuid "github.com/gofrs/uuid"
)

// OrderRepository is an autogenerated mock type for the OrderRepository type
type OrderRepository struct {
	mock.Mock
}

// AddOrderItemIntoOrder provides a mock function with given fields: ctx, userID, items
func (_m *OrderRepository) AddOrderItemIntoOrder(ctx context.Context, userID string, items *[]orderItem.OrderItem) error {
	ret := _m.Called(ctx, userID, items)

	if len(ret) == 0 {
		panic("no return value specified for AddOrderItemIntoOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *[]orderItem.OrderItem) error); ok {
		r0 = rf(ctx, userID, items)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AlterUserOrder provides a mock function with given fields: ctx, userID
func (_m *OrderRepository) AlterUserOrder(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for AlterUserOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AlterUserOrderByID provides a mock function with given fields: ctx, userID, orderID
func (_m *OrderRepository) AlterUserOrderByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) error {
	ret := _m.Called(ctx, userID, orderID)

	if len(ret) == 0 {
		panic("no return value specified for AlterUserOrderByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, userID, orderID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ChangeStatusOfCart provides a mock function with given fields: ctx, userId, orderId
func (_m *OrderRepository) ChangeStatusOfCart(ctx context.Context, userId string, orderId string) error {
	ret := _m.Called(ctx, userId, orderId)

	if len(ret) == 0 {
		panic("no return value specified for ChangeStatusOfCart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, userId, orderId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckOrderExists provides a mock function with given fields: ctx, userID
func (_m *OrderRepository) CheckOrderExists(ctx context.Context, userID string) (bool, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for CheckOrderExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUserOrder provides a mock function with given fields: ctx, userID
func (_m *OrderRepository) CreateUserOrder(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for CreateUserOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetOrderItemsFromOrderID provides a mock function with given fields: ctx, orderID
func (_m *OrderRepository) GetOrderItemsFromOrderID(ctx context.Context, orderID string) (*[]orderItem.OrderItemFull, error) {
	ret := _m.Called(ctx, orderID)

	if len(ret) == 0 {
		panic("no return value specified for GetOrderItemsFromOrderID")
	}

	var r0 *[]orderItem.OrderItemFull
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*[]orderItem.OrderItemFull, error)); ok {
		return rf(ctx, orderID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *[]orderItem.OrderItemFull); ok {
		r0 = rf(ctx, orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]orderItem.OrderItemFull)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrdersByUserID provides a mock function with given fields: ctx, userID
func (_m *OrderRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]orderItem.HistoryOrderItem, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetOrdersByUserID")
	}

	var r0 []orderItem.HistoryOrderItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]orderItem.HistoryOrderItem, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []orderItem.HistoryOrderItem); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]orderItem.HistoryOrderItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserOrderByUserID provides a mock function with given fields: ctx, orderId
func (_m *OrderRepository) GetUserOrderByUserID(ctx context.Context, orderId string) (*order.Model, error) {
	ret := _m.Called(ctx, orderId)

	if len(ret) == 0 {
		panic("no return value specified for GetUserOrderByUserID")
	}

	var r0 *order.Model
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*order.Model, error)); ok {
		return rf(ctx, orderId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *order.Model); ok {
		r0 = rf(ctx, orderId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*order.Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, orderId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsersOrder provides a mock function with given fields: ctx, userId
func (_m *OrderRepository) GetUsersOrder(ctx context.Context, userId string) (*order.Model, error) {
	ret := _m.Called(ctx, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetUsersOrder")
	}

	var r0 *order.Model
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*order.Model, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *order.Model); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*order.Model)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveCartItem provides a mock function with given fields: ctx, userID, bookID
func (_m *OrderRepository) RemoveCartItem(ctx context.Context, userID string, bookID uuid.UUID) error {
	ret := _m.Called(ctx, userID, bookID)

	if len(ret) == 0 {
		panic("no return value specified for RemoveCartItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uuid.UUID) error); ok {
		r0 = rf(ctx, userID, bookID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewOrderRepository creates a new instance of OrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderRepository {
	mock := &OrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
