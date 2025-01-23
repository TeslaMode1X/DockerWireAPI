package middleware

import (
	"context"
	"fmt"
	"github.com/TeslaMode1X/DockerWireAPI/internal/utils/response"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
	"os"
)

func WithAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := getTokenFromRequest(r)
		if err != nil {
			permissionDenied(w, r, "unable to get token from request")
		}

		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			permissionDenied(w, r, "invalid token")
		}

		userID, err := getUserIDFromToken(token)
		if err != nil {
			permissionDenied(w, r, "unable to get user ID from token")
		}

		// Store the user ID in the request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func permissionDenied(w http.ResponseWriter, r *http.Request, error string) {
	response.WriteError(w, r, http.StatusUnauthorized, slog.String("error", error))
	return
}

func getTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt-token")
	if err != nil {
		return "", err
	}

	tokenString := cookie.Value

	return tokenString, nil
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY_AUTH")), nil
	})
}

func getUserIDFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	return userID, nil
}
