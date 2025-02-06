package user

import (
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces/mocks"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/user"
	"github.com/TeslaMode1X/DockerWireAPI/internal/service"
	"github.com/TeslaMode1X/DockerWireAPI/packages/logger"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_NewUserHandler(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.UserService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()

	t.Run("it should return 401 Unauthorized", func(t *testing.T) {
		hdl.NewUserHandler(router)
	})

}

func TestHandler_GetUserByID_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.UserService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/{id}", hdl.GetUserByID)

	id, _ := uuid.NewV4()

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/"+id.String(), nil)

		var user user.User

		svc.On("GetUserByID", mock.Anything, id.String()).Return(user, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func TestHandler_GetUserByID_Svc_User_Not_Found_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.UserService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/{id}", hdl.GetUserByID)

	id, _ := uuid.NewV4()

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/"+id.String(), nil)

		var user user.User

		svc.On("GetUserByID", mock.Anything, id.String()).Return(user, service.ErrNotFound)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusNotFound, r.Code)
	})
}

func TestHandler_GetUserByID_Svc_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.UserService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/{id}", hdl.GetUserByID)

	id, _ := uuid.NewV4()

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/"+id.String(), nil)

		var user user.User

		svc.On("GetUserByID", mock.Anything, id.String()).Return(user, errors.New("user not found"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_GetUserByID_User_Id_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.UserService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/{id}", hdl.GetUserByID)

	id, _ := uuid.NewV4()

	t.Run("it should return 200 OK", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/"+"123", nil)

		var user user.User

		svc.On("GetUserByID", mock.Anything, id.String()).Return(user, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}
