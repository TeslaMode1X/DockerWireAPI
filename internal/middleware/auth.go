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
			return
		}

		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			permissionDenied(w, r, "invalid token")
			return
		}

		userID, err := getUserIDFromToken(token)
		if err != nil {
			permissionDenied(w, r, "unable to get user ID from token")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func WithOptionalAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := getTokenFromRequest(r)
		if err != nil {
			handler.ServeHTTP(w, r)
			return
		}

		token, err := validateToken(tokenString)
		if err != nil || !token.Valid {
			handler.ServeHTTP(w, r)
			return
		}

		userID, err := getUserIDFromToken(token)
		if err != nil {
			handler.ServeHTTP(w, r)
			return
		}

		role, err := getRoleFromToken(token)
		if err != nil {
			handler.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		ctx = context.WithValue(ctx, "role", role)

		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	})
}

func permissionDenied(w http.ResponseWriter, r *http.Request, error string) {
	response.WriteError(w, r, http.StatusUnauthorized, slog.String("error", error))
}

func getRoleFromToken(token *jwt.Token) (float64, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	role, ok := claims["role"].(float64)
	if !ok {
		return 0, errors.New("invalid user role in token")
	}
	return role, nil
}

func getTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt-token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
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

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

func GetUserRoleFromContext(ctx context.Context) (float64, error) {
	role, ok := ctx.Value("role").(float64)
	if !ok {
		return 0, errors.New("user role not found in context")
	}
	return role, nil
}
