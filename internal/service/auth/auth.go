package auth

import (
	"context"
	"fmt"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/auth"
	"github.com/TeslaMode1X/DockerWireAPI/internal/service"
	"github.com/gofrs/uuid"
)

type Service struct {
	AuthRepo interfaces.AuthRepository
	UserRepo interfaces.UserRepository
}

func (s *Service) Register(ctx context.Context, user model.Registration) (uuid.UUID, error) {
	const op = "service.user.Registration"

	userExists, err := s.UserRepo.CheckUserExists(ctx, user.Username)
	if err != nil {
		return uuid.Nil, err
	}
	if userExists {
		return uuid.Nil, fmt.Errorf("%s: %w", op, service.ErrValid)
	}

	id, err := s.AuthRepo.Register(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, service.ErrNotFound)
	}

	userFound, err := s.UserRepo.FindUserByID(ctx, id)
	if err != nil {
		return uuid.Nil, err
	}

	if userFound.ID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, service.ErrNotFound)
	}

	return id, err
}

func (s *Service) Login(ctx context.Context, user model.Login) (uuid.UUID, error) {
	const op = "service.user.Login"
	userExists, err := s.UserRepo.CheckUserExists(ctx, user.Username)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	if !userExists {
		return uuid.Nil, fmt.Errorf("%s: %w", op, service.ErrNotFound)
	}

	id, err := s.AuthRepo.Login(ctx, user)
	if err != nil {
		return id, fmt.Errorf("%s: %w", op, err)
	}

	return id, err
}
