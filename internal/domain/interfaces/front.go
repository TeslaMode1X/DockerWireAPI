package interfaces

import (
	"context"
	modelB "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
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
		ProcessLogin(ctx context.Context, w http.ResponseWriter, r *http.Request, form url.Values) error
		AdminPage(ctx context.Context, params mainPageParams.Model) (string, error)
		EditBook(ctx context.Context, bookID string, book *modelB.Book) error
		DeleteBook(ctx context.Context, bookID string) error
	}
)

type (
	FrontHandler interface {
		MainPage(w http.ResponseWriter, r *http.Request)
		RegistrationPage(w http.ResponseWriter, r *http.Request)
		RegistrationFront(w http.ResponseWriter, r *http.Request)
		LoginPage(w http.ResponseWriter, t *http.Request)
		LoginFront(w http.ResponseWriter, t *http.Request)
		AdminPage(w http.ResponseWriter, r *http.Request)
		EditBookFront(w http.ResponseWriter, r *http.Request)
		DeleteBookFront(w http.ResponseWriter, r *http.Request)
	}
)
