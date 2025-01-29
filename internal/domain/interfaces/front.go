package interfaces

import (
	"context"
	"net/http"
)

type (
	FrontRepository interface {
		MainPage(ctx context.Context)
	}
)

type (
	FrontService interface {
		MainPage(ctx context.Context, page string) (string, error)
	}
)

type (
	FrontHandler interface {
		MainPage(w http.ResponseWriter, r *http.Request)
	}
)
