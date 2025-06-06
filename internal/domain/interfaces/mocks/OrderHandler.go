// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// OrderHandler is an autogenerated mock type for the OrderHandler type
type OrderHandler struct {
	mock.Mock
}

// AddOrderItemIntoOrder provides a mock function with given fields: w, r
func (_m *OrderHandler) AddOrderItemIntoOrder(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// AlterUserOrder provides a mock function with given fields: w, r
func (_m *OrderHandler) AlterUserOrder(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// CreateUserOrder provides a mock function with given fields: w, r
func (_m *OrderHandler) CreateUserOrder(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetUserOrderByUserID provides a mock function with given fields: w, r
func (_m *OrderHandler) GetUserOrderByUserID(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// GetUsersOrder provides a mock function with given fields: w, r
func (_m *OrderHandler) GetUsersOrder(w http.ResponseWriter, r *http.Request) {
	_m.Called(w, r)
}

// NewOrderHandler creates a new instance of OrderHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderHandler {
	mock := &OrderHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
