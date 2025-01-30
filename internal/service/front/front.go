package front

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/auth"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

type Service struct {
	UserRepo  interfaces.UserRepository
	AuthRepo  interfaces.AuthRepository
	BookRepo  interfaces.BookRepository
	Templates map[string]*template.Template
}

func (s *Service) MainPage(ctx context.Context, page, errorMessage, successMessage string) (string, error) {
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
		"Title":   "Welcome to " + page,
		"Books":   allBooks,
		"Error":   errorMessage,
		"Success": successMessage,
	})
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	return buf.String(), nil
}

func (s *Service) LoginPage(ctx context.Context, page, errorMessage, successMessage string) (string, error) {
	const op = "service.front.MainPage"

	var tmpl, ok = s.Templates[page]
	if !ok {
		return "", errors.Wrap(errors.New("couldn't load template"), op)
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, map[string]interface{}{
		"Title":   "Welcome to " + page,
		"Error":   errorMessage,
		"Success": successMessage,
	})
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	return buf.String(), nil
}

func (s *Service) ProcessLogin(ctx context.Context, form url.Values) error {
	const op = "service.front.ProcessLogin"

	email := form.Get("email")
	password := form.Get("password")

	if email == "" || password == "" {
		return errors.New("empty_fields")
	}

	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return errors.New("json_encoding_failed")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("request_failed")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("server_unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid_credentials")
	}

	return nil
}

func (s *Service) RegistrationPage(ctx context.Context, page, errorMessage, successMessage string) (string, error) {
	const op = "service.front.MainPage"

	var tmpl, ok = s.Templates[page]
	if !ok {
		return "", errors.Wrap(errors.New("couldn't load template"), op)
	}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, map[string]interface{}{
		"Title":   "Welcome to " + page,
		"Error":   errorMessage,
		"Success": successMessage,
	})
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	return buf.String(), nil
}

func (s *Service) ProcessRegistration(ctx context.Context, form url.Values) error {
	const op = "service.front.ProcessRegistration"

	username := form.Get("username")
	email := form.Get("email")
	password := form.Get("password")

	if email == "" || password == "" {
		return errors.New("empty_fields")
	}

	if len(strings.TrimSpace(password)) < 3 {
		return errors.New("password is too short")
	}

	user := model.Registration{
		Username: username,
		Email:    email,
		Password: password,
	}

	hashedPassword := sha256.Sum256([]byte(user.Password))
	user.Password = hex.EncodeToString(hashedPassword[:])

	_, err := s.AuthRepo.Register(ctx, user)
	if err != nil {
		return errors.New("user by that mail already exists")
	}

	return nil
}
