package interfaces

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/user"
	"github.com/gofrs/uuid"
	"net/http"
)

type (
	UserRepository interface {
		CheckUserExists(ctx context.Context, username string) (bool, error)
		FindUserByID(ctx context.Context, id uuid.UUID) (user.User, error)
	}
)

type (
	UserService interface {
		GetUserByID(ctx context.Context, id string) (user.User, error)
	}
)

type (
	UserHandler interface {
		GetUserByID(w http.ResponseWriter, r *http.Request)
	}
)
