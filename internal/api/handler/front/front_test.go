package front

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces/mocks"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/mainPageParams"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"
	"github.com/TeslaMode1X/DockerWireAPI/packages/logger"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler_NewFrontEndHandler(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()

	t.Run("should be no errors", func(t *testing.T) {
		hdl.NewFrontEndHandler(router)
	})
}

func TestHandler_MainPage_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/", hdl.MainPage)

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		params := mainPageParams.Model{
			Page: "main",
		}
		svc.On("MainPage", mock.Anything, params).Return("<html>Test</html>", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Contains(t, r.Body.String(), "<html>Test</html>")
	})
}

func TestHandler_MainPage_InternalServerError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/", hdl.MainPage)

	t.Run("it should return 500 Internal Server Error", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		params := mainPageParams.Model{
			Page: "main",
		}
		svc.On("MainPage", mock.Anything, params).Return("", errors.New("internal error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_AddCartItems_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/cart/add", hdl.AddCartItems)

	t.Run("it should redirect with success message", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cart/add?id=123e4567-e89b-12d3-a456-426614174000&quantity=2", nil)

		ctx := context.WithValue(req.Context(), "user_id", "test_user_id")
		req = req.WithContext(ctx)

		items := []orderItem.OrderItem{
			{
				BookID:   uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000")),
				Quantity: 2,
			},
		}
		svc.On("AddCartItems", mock.Anything, "test_user_id", &items).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/?success=added_to_cart")
	})
}

func TestHandler_AddCartItems_InvalidBookID(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/cart/add", hdl.AddCartItems)

	t.Run("it should return 400 Bad Request for invalid book ID", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cart/add?id=invalid-id", nil)

		ctx := context.WithValue(req.Context(), "user_id", "test_user_id")
		req = req.WithContext(ctx)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_LoginFront_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/login/front", hdl.LoginFront)

	t.Run("it should redirect with success message", func(t *testing.T) {
		r := httptest.NewRecorder()
		form := url.Values{}
		form.Set("username", "test_user")
		form.Set("password", "test_password")
		req, _ := http.NewRequest(http.MethodPost, "/login/front", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		svc.On("ProcessLogin", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/?success=logged_in")
	})
}

func TestHandler_LoginFront_FormError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/login/front", hdl.LoginFront)

	t.Run("it should redirect with form error", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login/front", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/login?error=form_error")
	})
}

func TestHandler_LoginPage_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/login", hdl.LoginPage)

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/login?error=test_error&success=test_success", nil)

		svc.On("LoginPage", mock.Anything, "login", "test_error", "test_success").Return("<html>Login Page</html>", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Contains(t, r.Body.String(), "<html>Login Page</html>")
	})
}

func TestHandler_LoginPage_InternalServerError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/login", hdl.LoginPage)

	t.Run("it should return 500 Internal Server Error", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/login", nil)

		svc.On("LoginPage", mock.Anything, "login", "", "").Return("", errors.New("internal error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_RegistrationPage_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/register", hdl.RegistrationPage)

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/register?error=test_error&success=test_success", nil)

		svc.On("RegistrationPage", mock.Anything, "registration", "test_error", "test_success").Return("<html>Registration Page</html>", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Contains(t, r.Body.String(), "<html>Registration Page</html>")
	})
}

func TestHandler_RegistrationPage_InternalServerError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/register", hdl.RegistrationPage)

	t.Run("it should return 500 Internal Server Error", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/register", nil)

		svc.On("RegistrationPage", mock.Anything, "registration", "", "").Return("", errors.New("internal error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_RegistrationFront_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/register/front", hdl.RegistrationFront)

	t.Run("it should redirect with success message", func(t *testing.T) {
		r := httptest.NewRecorder()
		form := url.Values{}
		form.Set("username", "test_user")
		form.Set("password", "test_password")
		req, _ := http.NewRequest(http.MethodPost, "/register/front", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		svc.On("ProcessRegistration", mock.Anything, mock.Anything).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/login?success=registered")
	})
}

func TestHandler_RegistrationFront_FormError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/register/front", hdl.RegistrationFront)

	t.Run("it should return 400 Bad Request for form error", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register/front", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_AdminPage_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/admin", hdl.AdminPage)

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/admin?error=test_error&success=test_success", nil)

		params := mainPageParams.Model{
			Page:           "admin",
			ErrorMessage:   "test_error",
			SuccessMessage: "test_success",
		}
		svc.On("AdminPage", mock.Anything, params).Return("<html>Admin Page</html>", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Contains(t, r.Body.String(), "<html>Admin Page</html>")
	})
}

func TestHandler_AdminPage_InternalServerError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/admin", hdl.AdminPage)

	t.Run("it should return 500 Internal Server Error", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/admin", nil)

		params := mainPageParams.Model{
			Page: "admin",
		}
		svc.On("AdminPage", mock.Anything, params).Return("", errors.New("internal error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_HistoryPage_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/history", hdl.HistoryPage)

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/history", nil)

		ctx := context.WithValue(req.Context(), "user_id", "test_user_id")
		req = req.WithContext(ctx)

		svc.On("HistoryPage", mock.Anything, "test_user_id").Return("<html>History Page</html>", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Contains(t, r.Body.String(), "<html>History Page</html>")
	})
}

func TestHandler_HistoryPage_BadRequestError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/history", hdl.HistoryPage)

	t.Run("it should return 400 Bad Request for missing user_id", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/history", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_EditBookFront_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/admin/edit/{id}", hdl.EditBookFront)

	t.Run("it should redirect with success message", func(t *testing.T) {
		r := httptest.NewRecorder()
		form := url.Values{}
		form.Set("title", "New Title")
		form.Set("author", "New Author")
		form.Set("price", "19.99")
		form.Set("stock", "10")
		req, _ := http.NewRequest(http.MethodPost, "/admin/edit/123e4567-e89b-12d3-a456-426614174000", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		book := model.Book{
			Title:  "New Title",
			Author: "New Author",
			Price:  19.99,
			Stock:  10,
		}
		svc.On("EditBook", mock.Anything, "123e4567-e89b-12d3-a456-426614174000", &book).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/admin?success=book_updated")
	})
}

func TestHandler_DeleteBookFront_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/admin/delete/{id}", hdl.DeleteBookFront)

	t.Run("it should redirect with success message", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/admin/delete/123e4567-e89b-12d3-a456-426614174000", nil)

		svc.On("DeleteBook", mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/admin?success=book_deleted")
	})
}

func TestHandler_GetCartItems_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/api/v1/cart/items", hdl.GetCartItems)

	t.Run("it should return 200 OK with cart items", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/cart/items", nil)

		ctx := context.WithValue(req.Context(), "user_id", "test_user_id")
		req = req.WithContext(ctx)

		cartItems := []orderItem.OrderItemFull{
			{BookID: uuid.Must(uuid.NewV4()), Quantity: 2},
		}
		svc.On("GetCartItems", mock.Anything, "test_user_id").Return(&cartItems, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
		var response []orderItem.OrderItem
		err := json.Unmarshal(r.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, cartItems[0].Quantity, response[0].Quantity)
	})
}

func TestHandler_GetCartItems_Unauthorized(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/api/v1/cart/items", hdl.GetCartItems)

	t.Run("it should return 401 Unauthorized for missing user_id", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/cart/items", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})
}

func TestHandler_RemoveCartItem_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/api/v1/cart/remove", hdl.RemoveCartItem)

	t.Run("it should redirect with success message", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/cart/remove?id=123e4567-e89b-12d3-a456-426614174000", nil)

		ctx := context.WithValue(req.Context(), "user_id", "test_user_id")
		req = req.WithContext(ctx)

		bookID := uuid.Must(uuid.FromString("123e4567-e89b-12d3-a456-426614174000"))
		svc.On("RemoveCartItem", mock.Anything, "test_user_id", bookID).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/?success=removed_from_cart")
	})
}

func TestHandler_RemoveCartItem_InvalidBookID(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/api/v1/cart/remove", hdl.RemoveCartItem)

	t.Run("it should return 400 Bad Request for invalid book ID", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/cart/remove?id=invalid-id", nil)

		ctx := context.WithValue(req.Context(), "user_id", "test_user_id")
		req = req.WithContext(ctx)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_CartCheckout_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/cart/success", hdl.CartCheckout)

	t.Run("it should redirect with success message", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cart/success", nil)

		ctx := context.WithValue(req.Context(), "user_id", "test_user_id")
		req = req.WithContext(ctx)

		svc.On("CartCheckout", mock.Anything, "test_user_id").Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/?success=cart_paid_successfully")
	})
}

func TestHandler_CartCheckout_UserNotLoggedIn(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.FrontService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/cart/success", hdl.CartCheckout)

	t.Run("it should redirect to login for missing user_id", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/cart/success", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusSeeOther, r.Code)
		assert.Contains(t, r.Header().Get("Location"), "/login?error=user_not_logged_in")
	})
}
