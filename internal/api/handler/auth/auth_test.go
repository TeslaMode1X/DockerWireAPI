package auth

import (
	"encoding/json"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces/mocks"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/auth"
	"github.com/TeslaMode1X/DockerWireAPI/internal/service"
	"github.com/TeslaMode1X/DockerWireAPI/packages/logger"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_NewAuthHandler(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()

	t.Run("should be no errors", func(t *testing.T) {
		hdl.NewAuthHandler(router)
	})
}

func TestHandler_Registration_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/registration", hdl.Register)

	t.Run("should be no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		payload := `{"username": "test user", "email": "testuser@example.com", "password": "password123"}`

		req, _ := http.NewRequest("POST", "/registration", strings.NewReader(payload))

		var registration model.Registration
		err := json.Unmarshal([]byte(payload), &registration)
		assert.NoError(t, err)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), 12)
		assert.NoError(t, err)
		registration.Password = string(hashedPassword)

		svc.On("Register", mock.Anything, mock.Anything).Return(nil, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusCreated, r.Code)
	})
}

func TestHandler_Login_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/login", hdl.Login)

	t.Run("should return success on valid login", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/login", strings.NewReader(payload))
		require.NoError(t, err)

		userID, _ := uuid.NewV4()
		svc.On("Login", mock.Anything, mock.Anything).Return(userID, 1, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func TestHandler_Login_Model_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/login", hdl.Login)

	t.Run("should return success on valid login", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"emil": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/login", strings.NewReader(payload))
		require.NoError(t, err)

		userID, _ := uuid.NewV4()
		svc.On("Login", mock.Anything, mock.Anything).Return(userID, 1, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_Login_SVC_ErrNotFound(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/login", hdl.Login)

	t.Run("should return success on valid login", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/login", strings.NewReader(payload))
		require.NoError(t, err)

		userID, _ := uuid.NewV4()
		svc.On("Login", mock.Anything, mock.Anything).Return(userID, 1, service.ErrNotFound)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}

func TestHandler_Login_SVC_ErrValid(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/login", hdl.Login)

	t.Run("should return success on valid login", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/login", strings.NewReader(payload))
		require.NoError(t, err)

		userID, _ := uuid.NewV4()
		svc.On("Login", mock.Anything, mock.Anything).Return(userID, 1, service.ErrValid)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_Login_SVC_Err(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/login", hdl.Login)

	t.Run("should return success on valid login", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/login", strings.NewReader(payload))
		require.NoError(t, err)

		userID, _ := uuid.NewV4()
		svc.On("Login", mock.Anything, mock.Anything).Return(userID, 1, errors.New("test error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_Registration_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/registration", hdl.Register)

	t.Run("should return internal server error", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"username": "testuser", "email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/registration", strings.NewReader(payload))
		require.NoError(t, err)

		var registration model.Registration
		err = json.Unmarshal([]byte(payload), &registration)
		require.NoError(t, err)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), 12)
		require.NoError(t, err)
		registration.Password = string(hashedPassword)

		svc.On("Register", mock.Anything, mock.Anything).Return(nil, errors.New("user already exists"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
		body := r.Body.String()
		assert.Contains(t, body, service.ErrUserAlreadyExists.Error())
	})

	t.Run("should return user already exists", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"username": "testuser", "email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/registration", strings.NewReader(payload))
		require.NoError(t, err)

		var registration model.Registration
		err = json.Unmarshal([]byte(payload), &registration)
		require.NoError(t, err)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), 12)
		require.NoError(t, err)
		registration.Password = string(hashedPassword)

		svc.On("Register", mock.Anything, mock.Anything).Return(nil, service.ErrUserAlreadyExists)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
		body := r.Body.String()
		assert.Contains(t, body, service.ErrUserAlreadyExists.Error())
	})

	t.Run("should be model error", func(t *testing.T) {
		r := httptest.NewRecorder()

		payload := `{"name": "test user", "email": "testuser@example.com", "password": "password123"}`

		req, _ := http.NewRequest("POST", "/registration", strings.NewReader(payload))

		var registration model.Registration
		err := json.Unmarshal([]byte(payload), &registration)
		assert.NoError(t, err)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), 12)
		assert.NoError(t, err)
		registration.Password = string(hashedPassword)

		svc.On("Register", mock.Anything, mock.Anything).Return(nil, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})

}

func TestHandler_Registration_Svc_Error_User_Already_Exists(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/registration", hdl.Register)

	t.Run("should return user already exists", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"username": "testuser", "email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/registration", strings.NewReader(payload))
		require.NoError(t, err)

		var registration model.Registration
		err = json.Unmarshal([]byte(payload), &registration)
		require.NoError(t, err)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), 12)
		require.NoError(t, err)
		registration.Password = string(hashedPassword)

		svc.On("Register", mock.Anything, mock.Anything).Return(nil, service.ErrUserAlreadyExists)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
		body := r.Body.String()
		assert.Contains(t, body, service.ErrUserAlreadyExists.Error())
	})

}

func TestHandler_Registration_Svc_Error_Not_Found(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/registration", hdl.Register)

	t.Run("should return user not found", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"username": "testuser", "email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/registration", strings.NewReader(payload))
		require.NoError(t, err)

		var registration model.Registration
		err = json.Unmarshal([]byte(payload), &registration)
		require.NoError(t, err)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), 12)
		require.NoError(t, err)
		registration.Password = string(hashedPassword)

		svc.On("Register", mock.Anything, mock.Anything).Return(nil, service.ErrNotFound)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})

}

func TestHandler_Login_AlreadyLoggedIn(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.AuthService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Post("/login", hdl.Login)

	t.Run("should return 401 if user already logged in", func(t *testing.T) {
		r := httptest.NewRecorder()
		payload := `{"email": "testuser@example.com", "password": "password123"}`
		req, err := http.NewRequest("POST", "/login", strings.NewReader(payload))
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{
			Name:  "jwt-token",
			Value: "some-valid-token",
		})

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusUnauthorized, r.Code)
		assert.Contains(t, r.Body.String(), "already logged in")
	})
}
