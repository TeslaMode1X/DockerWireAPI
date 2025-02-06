package interfaces

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/user"
	"github.com/gofrs/uuid"
	"net/http"
)

//go:generate mockery --name UserRepository
type (
	UserRepository interface {
		CheckUserExists(ctx context.Context, username string) (bool, error)
		FindUserByID(ctx context.Context, id uuid.UUID) (user.User, error)
	}
)

//go:generate mockery --name UserService
type (
	UserService interface {
		GetUserByID(ctx context.Context, id string) (user.User, error)
	}
)

//go:generate mockery --name UserHandler
type (
	UserHandler interface {
		GetUserByID(w http.ResponseWriter, r *http.Request)
	}
)
