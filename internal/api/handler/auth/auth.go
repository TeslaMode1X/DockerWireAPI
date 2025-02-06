package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/TeslaMode1X/DockerWireAPI/internal/domain/interfaces"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/auth"
	"github.com/TeslaMode1X/DockerWireAPI/internal/service"
	"github.com/TeslaMode1X/DockerWireAPI/internal/utils/response"
	"github.com/TeslaMode1X/DockerWireAPI/packages/jsonReader"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"time"
)

type Handler struct {
	Svc interfaces.AuthService
	Log *slog.Logger
}

func (h *Handler) NewAuthHandler(r chi.Router) {
	r.Post("/registration", h.Register)
	r.Post("/login", h.Login)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.Registration"
	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var user model.Registration
	if err := jsonReader.ReadJSON(w, r, &user); err != nil {
		h.Log.Error("failed to decode json body", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, errors.New("failed to decode request body"))
		return
	}

	hashedPassword := sha256.Sum256([]byte(user.Password))
	user.Password = hex.EncodeToString(hashedPassword[:])

	userCreated, err := h.Svc.Register(context.Background(), user)
	if err != nil {
		h.Log.Error("Error during registration", slog.String("error", err.Error()))

		if errors.Is(err, service.ErrUserAlreadyExists) {
			response.WriteError(w, r, http.StatusBadRequest, fmt.Sprintf("%v", service.ErrUserAlreadyExists))
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			response.WriteError(w, r, http.StatusBadRequest, fmt.Sprintf("%v", service.ErrNotFound))
			return
		}
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	response.WriteJson(w, r, http.StatusCreated, userCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.Login"

	h.Log = h.Log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	cookie, err := r.Cookie("jwt-token")
	if err == nil {
		response.WriteError(w, r, http.StatusUnauthorized, errors.New("already logged in"))
		return
	}

	var userCurrent model.Login
	if err := jsonReader.ReadJSON(w, r, &userCurrent); err != nil {
		h.Log.Error("failed to decode request body", slog.String("error", err.Error()))
		response.WriteError(w, r, http.StatusBadRequest, errors.New("failed to decode request body"))
		return
	}

	userID, role, err := h.Svc.Login(context.Background(), userCurrent)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.WriteError(w, r, http.StatusNotFound, service.ErrNotFound)
			return
		}
		if errors.Is(err, service.ErrValid) {
			response.WriteError(w, r, http.StatusBadRequest, service.ErrValid)
			return
		}
		response.WriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // token expires in 1 hour
		"user_id": userID,
		"role":    role,
	})
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	cookie = &http.Cookie{
		Name:     "jwt-token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 1),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	response.WriteJson(w, r, http.StatusOK, userID.String())
}
