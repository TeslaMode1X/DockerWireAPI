package front

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/auth"
	modelB "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/mainPageParams"
	orderModels "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type Service struct {
	UserRepo  interfaces.UserRepository
	AuthRepo  interfaces.AuthRepository
	BookRepo  interfaces.BookRepository
	OrderRepo interfaces.OrderRepository
	Templates map[string]*template.Template
}

func (s *Service) MainPage(ctx context.Context, params mainPageParams.Model) (string, error) {
	const op = "service.front.MainPage"

	tmpl, ok := s.Templates[params.Page]
	if !ok {
		return "", errors.Wrap(errors.New("couldn't load template"), op)
	}

	allBooks, err := s.BookRepo.GetAllBooks(ctx)
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	var filteredBooks []modelB.Book
	if params.SearchQuery != "" {
		for _, book := range *allBooks {
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(params.SearchQuery)) {
				filteredBooks = append(filteredBooks, book)
			}
		}
	} else {
		filteredBooks = *allBooks
	}

	switch params.SortBy {
	case "name":
		sort.Slice(filteredBooks, func(i, j int) bool {
			return strings.ToLower(filteredBooks[i].Title) < strings.ToLower(filteredBooks[j].Title)
		})
	case "price":
		sort.Slice(filteredBooks, func(i, j int) bool {
			return filteredBooks[i].Price < filteredBooks[j].Price
		})
	case "stock":
		sort.Slice(filteredBooks, func(i, j int) bool {
			return filteredBooks[i].Stock < filteredBooks[j].Stock
		})
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Title":       "Welcome to " + params.Page,
		"Books":       filteredBooks,
		"Error":       params.ErrorMessage,
		"Success":     params.SuccessMessage,
		"SearchQuery": params.SearchQuery,
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

func (s *Service) ProcessLogin(ctx context.Context, w http.ResponseWriter, r *http.Request, form url.Values) error {
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

	for _, cookie := range r.Cookies() {
		req.AddCookie(cookie)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("server_unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("already_logged_in")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid_credentials")
	}

	for _, cookie := range resp.Cookies() {
		http.SetCookie(w, cookie)
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

func (s *Service) AdminPage(ctx context.Context, params mainPageParams.Model) (string, error) {
	const op = "service.front.AdminPage"

	tmpl, ok := s.Templates[params.Page]
	if !ok {
		return "", errors.Wrap(errors.New("couldn't load template"), op)
	}

	allBooks, err := s.BookRepo.GetAllBooks(ctx)
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	var filteredBooks []modelB.Book
	if params.SearchQuery != "" {
		for _, book := range *allBooks {
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(params.SearchQuery)) {
				filteredBooks = append(filteredBooks, book)
			}
		}
	} else {
		filteredBooks = *allBooks
	}

	switch params.SortBy {
	case "name":
		sort.Slice(filteredBooks, func(i, j int) bool {
			return strings.ToLower(filteredBooks[i].Title) < strings.ToLower(filteredBooks[j].Title)
		})
	case "price":
		sort.Slice(filteredBooks, func(i, j int) bool {
			return filteredBooks[i].Price < filteredBooks[j].Price
		})
	case "stock":
		sort.Slice(filteredBooks, func(i, j int) bool {
			return filteredBooks[i].Stock < filteredBooks[j].Stock
		})
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Title":       "Admin Panel",
		"Books":       filteredBooks,
		"Error":       params.ErrorMessage,
		"Success":     params.SuccessMessage,
		"SearchQuery": params.SearchQuery,
	})
	if err != nil {
		return "", errors.Wrap(err, op)
	}

	return buf.String(), nil
}

func (s *Service) EditBook(ctx context.Context, bookID string, book *modelB.Book) error {
	const op = "service.front.EditBook"

	jsonData, err := json.Marshal(book)
	if err != nil {
		return errors.Wrap(err, op)
	}

	req, err := http.NewRequest("PUT", "http://localhost:8080/api/v1/book/"+bookID, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, op)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("failed to update book"), op)
	}

	return nil
}

func (s *Service) DeleteBook(ctx context.Context, bookID string) error {
	const op = "service.front.DeleteBook"

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/book/"+bookID, nil)
	if err != nil {
		return errors.Wrap(err, op)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, op)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(errors.New("failed to delete book"), op)
	}

	return nil
}

func (s *Service) GetCartItems(ctx context.Context, userId string) (*[]orderModels.OrderItemFull, error) {
	const op = "service.front.GetCartItems"

	exists, err := s.OrderRepo.CheckOrderExists(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	if !exists {
		err = s.OrderRepo.CreateUserOrder(ctx, userId)
		fmt.Println(err)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}
	}

	userOrder, err := s.OrderRepo.GetUsersOrder(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	orderID := userOrder.ID.String()

	orderItems, err := s.OrderRepo.GetOrderItemsFromOrderID(ctx, orderID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return orderItems, nil
}

func (s *Service) AddCartItems(ctx context.Context, userID string, items *[]orderModels.OrderItem) error {
	const op = "service.front.AddCartItems"

	err := s.OrderRepo.AddOrderItemIntoOrder(ctx, userID, items)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (s *Service) RemoveCartItem(ctx context.Context, userID string, bookID uuid.UUID) error {
	const op = "service.front.RemoveCartItem"

	err := s.OrderRepo.RemoveCartItem(ctx, userID, bookID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
