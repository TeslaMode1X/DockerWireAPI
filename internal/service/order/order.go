package order

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	orderModel "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/order"
	orderModels "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"
	"github.com/pkg/errors"
)

type Service struct {
	OrderRepo interfaces.OrderRepository
}

func (s *Service) GetUsersOrder(ctx context.Context, userId string) (*orderModel.Model, error) {
	const op = "service.order.GetUsersOrder"

	order, err := s.OrderRepo.GetUsersOrder(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return order, nil
}

func (s *Service) GetUserOrderByUserID(ctx context.Context, orderId string) (*orderModel.Model, error) {
	const op = "service.order.GetUsersOrder"

	order, err := s.OrderRepo.GetUserOrderByUserID(ctx, orderId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return order, nil
}

func (s *Service) CreateUserOrder(ctx context.Context, userID string) error {
	const op = "service.order.CreateUserOrder"

	exists, err := s.OrderRepo.CheckOrderExists(ctx, userID)
	if err != nil {
		return errors.Wrap(err, op)
	}
	if exists {
		return errors.Wrap(errors.New("order already exists"), op)
	}

	err = s.OrderRepo.CreateUserOrder(ctx, userID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s *Service) AlterUserOrder(ctx context.Context, userID string) error {
	const op = "service.order.AlterUserOrder"

	exists, err := s.OrderRepo.CheckOrderExists(ctx, userID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	if !exists {
		return errors.Wrap(errors.New("order does not exists"), op)
	}

	err = s.OrderRepo.AlterUserOrder(ctx, userID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s *Service) AddOrderItemIntoOrder(ctx context.Context, userID string, items *[]orderModels.OrderItem) error {
	const op = "service.order.AddOrderItemIntoOrder"

	err := s.OrderRepo.AddOrderItemIntoOrder(ctx, userID, items)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
