package interfaces

import (
	"context"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/auth"
	"github.com/gofrs/uuid"
	"net/http"
)

//go:generate mockery --name AuthRepository
type (
	AuthRepository interface {
		Register(ctx context.Context, user model.Registration) (uuid.UUID, error)
		Login(ctx context.Context, user model.Login) (uuid.UUID, int, error)
	}
)

//go:generate mockery --name AuthService
type (
	AuthService interface {
		Register(ctx context.Context, user model.Registration) (uuid.UUID, error)
		Login(ctx context.Context, user model.Login) (uuid.UUID, int, error)
	}
)

//go:generate mockery --name AuthHandler
type (
	AuthHandler interface {
		Login(w http.ResponseWriter, r *http.Request)
		Register(w http.ResponseWriter, r *http.Request)
	}
)
