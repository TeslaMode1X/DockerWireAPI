package interfaces

import (
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/mainPageParams"
	"net/http"
	"net/url"
)

type (
	FrontRepository interface {
		MainPage(ctx context.Context)
	}
)

type (
	FrontService interface {
		MainPage(ctx context.Context, model mainPageParams.Model) (string, error)
		RegistrationPage(ctx context.Context, page, errorMessage, successMessage string) (string, error)
		LoginPage(ctx context.Context, page, errorMessage, successMessage string) (string, error)
		ProcessRegistration(ctx context.Context, form url.Values) error
		ProcessLogin(ctx context.Context, form url.Values) error
	}
)

type (
	FrontHandler interface {
		MainPage(w http.ResponseWriter, r *http.Request)
		RegistrationPage(w http.ResponseWriter, r *http.Request)
		RegistrationFront(w http.ResponseWriter, r *http.Request)
		LoginPage(w http.ResponseWriter, t *http.Request)
		LoginFront(w http.ResponseWriter, t *http.Request)
	}
)
