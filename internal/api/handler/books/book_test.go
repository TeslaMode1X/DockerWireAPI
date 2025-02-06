package books

import (
	"bytes"
	"encoding/json"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces/mocks"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/books"
	"github.com/TeslaMode1X/DockerWireAPI/packages/logger"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_NewBookHandler(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()

	t.Run("it should return no errors", func(t *testing.T) {
		hdl.NewBookHandler(router)
	})

}

func TestHandler_GetBookHandler_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/", hdl.GetAllBooks)

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		svc.On("GetAllBooks", mock.Anything).Return(nil, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func TestHandler_GetBookHandler_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/", hdl.GetAllBooks)

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		svc.On("GetAllBooks", mock.Anything).Return(nil, errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_GetBookHandlerById_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/{id}", hdl.GetBookById)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, "/"+id.String(), nil)
		require.NoError(t, err)

		mockBook := &model.Book{
			ID:     id,
			Title:  "Test Book",
			Author: "Test Author",
		}
		svc.On("GetBookById", mock.Anything, id).Return(mockBook, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)

		body := r.Body.String()
		assert.Contains(t, body, `"id":"`+id.String()+`"`)
		assert.Contains(t, body, `"title":"Test Book"`)
		assert.Contains(t, body, `"author":"Test Author"`)
	})
}

func TestHandler_GetBookHandlerById_SVC_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Get("/{id}", hdl.GetBookById)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, "/"+id.String(), nil)
		require.NoError(t, err)

		svc.On("GetBookById", mock.Anything, id).Return(nil, errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)

	})
}

func TestHandler_GetBookHandlerById_UUID_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Get("/{id}", hdl.GetBookById)

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodGet, "/"+"123", nil)

		svc.On("GetBookById", mock.Anything, mock.Anything).Return(nil, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_CreateBook_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/", hdl.CreateBook)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		mockBook := model.Book{
			ID:     id,
			Title:  "Test Book",
			Author: "Test Author",
		}

		payload, err := json.Marshal(mockBook)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		svc.On("CreateBook", mock.Anything, mockBook).Return(id, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusCreated, r.Code)
	})
}

func TestHandler_CreateBook_SVC_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/", hdl.CreateBook)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		mockBook := model.Book{
			ID:     id,
			Title:  "Test Book",
			Author: "Test Author",
		}

		payload, err := json.Marshal(mockBook)
		require.NoError(t, err, "failed to marshal book data to JSON")

		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		svc.On("CreateBook", mock.Anything, mockBook).Return(id, errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_CreateBook_JsonMarshall_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Post("/", hdl.CreateBook)

	t.Run("should return 500 if response marshalling fails", func(t *testing.T) {
		r := httptest.NewRecorder()

		type BadBook struct {
			Title  string
			Author string
			Parent *BadBook `json:"parent"`
		}
		badBook := &BadBook{
			Title:  "Invalid Book",
			Author: "Author",
		}
		badBook.Parent = badBook

		payload, err := json.Marshal(badBook)
		require.Error(t, err, "unexpected error in test setup")

		req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		svc.On("CreateBook", mock.Anything, mock.Anything).Return(nil, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_UpdateBook_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Put("/{id}", hdl.UpdateBookById)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		payload := `{"title": "Updated Book", "author": "Updated Author"}`
		req, err := http.NewRequest(http.MethodPut, "/"+id.String(), strings.NewReader(payload))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		svc.On("UpdateBookById", mock.Anything, mock.Anything, mock.Anything).Return(id, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func TestHandler_UpdateBook_IdError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Put("/{id}", hdl.UpdateBookById)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		payload := `{"title": "Updated Book", "author": "Updated Author"}`
		req, err := http.NewRequest(http.MethodPut, "/"+"123", strings.NewReader(payload))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		svc.On("UpdateBookById", mock.Anything, mock.Anything, mock.Anything).Return(id, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_UpdateBook_JsonError(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Put("/{id}", hdl.UpdateBookById)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		payload := ``
		req, err := http.NewRequest(http.MethodPut, "/"+id.String(), strings.NewReader(payload))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		svc.On("UpdateBookById", mock.Anything, mock.Anything, mock.Anything).Return(id, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}

func TestHandler_UpdateBook_Svc_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Put("/{id}", hdl.UpdateBookById)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		payload := `{"title": "Updated Book", "author": "Updated Author"}`
		req, err := http.NewRequest(http.MethodPut, "/"+id.String(), strings.NewReader(payload))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		svc.On("UpdateBookById", mock.Anything, mock.Anything, mock.Anything).Return(id, errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_DeleteBookById_Success(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Delete("/{id}", hdl.DeleteBookById)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodDelete, "/"+id.String(), nil)
		require.NoError(t, err)

		svc.On("DeleteBookById", mock.Anything, id).Return(nil, nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusOK, r.Code)
	})
}

func TestHandler_DeleteBookById_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}
	router := chi.NewRouter()
	router.Delete("/{id}", hdl.DeleteBookById)

	id, _ := uuid.NewV4()

	t.Run("it should return no errors", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodDelete, "/"+id.String(), nil)
		require.NoError(t, err)

		svc.On("DeleteBookById", mock.Anything, id).Return(errors.New("error"))

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusInternalServerError, r.Code)
	})
}

func TestHandler_DeleteBookById_UUID_Error(t *testing.T) {
	log := logger.New(logger.EnvLocal)
	svc := mocks.BookService{}
	hdl := Handler{
		Svc: &svc,
		Log: log,
	}

	router := chi.NewRouter()
	router.Delete("/{id}", hdl.DeleteBookById)

	t.Run("it should return id error", func(t *testing.T) {
		r := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodDelete, "/"+"123", nil)

		svc.On("DeleteBookById", mock.Anything, mock.Anything).Return(nil)

		router.ServeHTTP(r, req)

		assert.Equal(t, http.StatusBadRequest, r.Code)
	})
}
