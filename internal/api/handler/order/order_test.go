package order

import (
	"context"
	"errors"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces/mocks"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/order"
	"github.com/TeslaMode1X/DockerWireAPI/packages/logger"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_NewOrderHandler(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}

	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()

	t.Run("it should return no errors", func(t *testing.T) {
		hdl.NewOrderHandler(router)
	})
}

func TestHandler_GetUsersOrder_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/orders", hdl.GetUsersOrder)

	id, _ := uuid.NewV4()

	t.Run("error - no user ID in context", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/orders", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})

	t.Run("success - should return user's order", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/orders", nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		mockOrder := model.Model{ID: id, UserID: id}
		svc.On("GetUsersOrder", mock.Anything, "123").Return(&mockOrder, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
	})

}

func TestHandler_GetUsersOrder_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/orders", hdl.GetUsersOrder)

	id, _ := uuid.NewV4()

	t.Run("error - no user ID in context", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/orders", nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})

	t.Run("success - should return user's order", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/orders", nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		mockOrder := model.Model{ID: id, UserID: id}
		svc.On("GetUsersOrder", mock.Anything, "123").Return(&mockOrder, errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})

}

func TestHandler_GetUsersOrderByID_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/{orderId}", hdl.GetUserOrderByUserID)

	id, _ := uuid.NewV4()

	t.Run("success - should return user's order", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+id.String(), nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		mockOrder := model.Model{ID: id, UserID: id}

		svc.On("GetUserOrderByUserID", mock.Anything, mock.Anything).Return(&mockOrder, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func TestHandler_GetUsersOrderByID_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/{orderId}", hdl.GetUserOrderByUserID)

	id, _ := uuid.NewV4()

	t.Run("success - should return user's order", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/"+id.String(), nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("GetUserOrderByUserID", mock.Anything, mock.Anything).Return(nil, errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_CreateOrder_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()

	router.Post("/", hdl.CreateUserOrder)

	t.Run("should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/", nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("CreateUserOrder", mock.Anything, "123").Return(nil, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusCreated, r.Code)
	})
}

func TestHandler_CreateOrder_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/", hdl.CreateUserOrder)

	t.Run("should return an error", func(t *testing.T) {
		r := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/", nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("CreateUserOrder", mock.Anything, "123").
			Return(errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_CreateOrder_User_Id_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()

	router.Post("/", hdl.CreateUserOrder)

	t.Run("error - no user ID in context", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/", nil)

		svc.On("CreateUserOrder", mock.Anything, "123").Return(nil, errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})
}

func TestHandler_AlterUserOrder_User_Id_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()

	router.Post("/", hdl.AlterUserOrder)

	t.Run("error - no user ID in context", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/", nil)

		svc.On("AlterUserOrder", mock.Anything, "123").Return(nil, errors.New("error"))
		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})
}

func TestHandler_AlterUserOrder_Object_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()

	router.Post("/{orderID}", hdl.AlterUserOrder)

	id, _ := uuid.NewV4()

	t.Run("error - no order ID in context", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/"+id.String(), nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("AlterUserOrder", mock.Anything, "123").Return(nil, errors.New("error"))
		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_AlterUserOrder_Svc_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/{orderId}", hdl.AlterUserOrder)

	id, _ := uuid.NewV4()

	t.Run("should return success with valid orderId", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/"+id.String(), nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("AlterUserOrderByID", mock.Anything, "123", id.String()).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func TestHandler_AlterUserOrder_Svc_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/{orderId}", hdl.AlterUserOrder)

	id, _ := uuid.NewV4()

	t.Run("should return success with valid orderId", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/"+id.String(), nil)

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("AlterUserOrderByID", mock.Anything, "123", id.String()).Return(errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_AddOrderItemIntoOrder_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/", hdl.AddOrderItemIntoOrder)

	t.Run("should return success with valid userId and valid request body", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"items": [{"product_id": "123", "quantity": 2}]}`))

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("AddOrderItemIntoOrder", mock.Anything, "123", mock.Anything).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusCreated, r.Code)
	})
}

func TestHandler_AddOrderItemIntoOrder_NoUserID(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/", hdl.AddOrderItemIntoOrder)

	t.Run("should return unauthorized if user_id is missing", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/", nil)

		svc.On("AddOrderItemIntoOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusUnauthorized, r.Code)
	})
}

func TestHandler_AddOrderItemIntoOrder_NoPayload(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/", hdl.AddOrderItemIntoOrder)

	t.Run("should return bad request if payload is missing", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("AddOrderItemIntoOrder", mock.Anything, "123", mock.Anything).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_AddOrderItemIntoOrder_MissingItems(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.OrderService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/", hdl.AddOrderItemIntoOrder)

	t.Run("should return bad request if items field is missing", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))

		ctx := context.WithValue(req.Context(), "user_id", "123")
		req = req.WithContext(ctx)

		svc.On("AddOrderItemIntoOrder", mock.Anything, "123", mock.Anything).Return(errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}
