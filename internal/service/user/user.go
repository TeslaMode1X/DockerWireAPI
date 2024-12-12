package user

import (
	"context"
	"fmt"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/user"
	"github.com/gofrs/uuid"
)

type Service struct {
	Repo interfaces.UserRepository
}

func (s *Service) GetUserByID(ctx context.Context, idStr string) (user.User, error) {
	const op = "service.user.GetUserByID"

	id, err := uuid.FromString(idStr)
	if err != nil {
		return user.User{}, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.Repo.FindUserByID(ctx, id)
	if err != nil {
		return user.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}
