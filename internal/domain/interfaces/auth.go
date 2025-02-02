package interfaces

import (
	"context"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/auth"
	"github.com/gofrs/uuid"
	"net/http"
)

type (
	AuthRepository interface {
		Register(ctx context.Context, user model.Registration) (uuid.UUID, error)
		Login(ctx context.Context, user model.Login) (uuid.UUID, int, error)
	}
)

type (
	AuthService interface {
		Register(ctx context.Context, user model.Registration) (uuid.UUID, error)
		Login(ctx context.Context, user model.Login) (uuid.UUID, int, error)
	}
)

type (
	AuthHandler interface {
		Login(w http.ResponseWriter, r *http.Request)
		Register(w http.ResponseWriter, r *http.Request)
	}
)
