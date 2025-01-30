package front

import (
	"bytes"
	"context"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	"github.com/pkg/errors"
	"html/template"
)

type Service struct {
	UserRepo  interfaces.UserRepository
	BookRepo  interfaces.BookRepository
	Templates map[string]*template.Template
}

func (s *Service) MainPage(ctx context.Context, page string) (string, error) {
	const op = "service.front.MainPage"

	var tmpl, ok = s.Templates[page]
	if !ok {
		return "", errors.Wrap(errors.New("couldn't load template"), op)
	}

	allBooks, err := s.BookRepo.GetAllBooks(ctx)
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Title": "Welcome to " + page,
		"Books": allBooks,
	})
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	return buf.String(), nil
}
