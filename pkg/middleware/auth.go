package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type contextKey string

const UserIDKey contextKey = "userID"

// NewAuthMiddleware создает middleware для валидации JWT-токена.
// Секрет должен быть передан параметром (или через env/config), чтобы обеспечить расширяемость и тестируемость.
func NewAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: missing bearer token", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			userID, err := validateJWT(tokenStr, secret)
			if err != nil {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// userID передаётся дальше через context.Context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// validateJWT разбирает токен и проверяет подпись, возвращает userID, если токен валиден.
func validateJWT(tokenStr, secret string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// В данном примере ожидаем, что user_id числовой
		UUID, err := claims.GetSubject()
		if err != nil {
			return "", errors.New("user_id missing in token")
		}
		return UUID, nil
	}
	return "", errors.New("invalid token")
}

func ValidateUser(r *http.Request) (string, error) {
	var (
		userID string
		ok     bool
	)
	if userID, ok = r.Context().Value(UserIDKey).(string); !ok {
		// http.Error(w, "Invalid user", http.StatusUnauthorized)
		return "", fmt.Errorf("not authenticated user")
	}

	return userID, nil
}

func ValidateAccount(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	accountIDStr := vars["id"]
	if accountIDStr == "" {
		// http.Error(w, "could not get account", http.StatusInternalServerError)
		return 0, fmt.Errorf(" missing account id")
	}
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		// http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return 0, fmt.Errorf(" missing account id")
	}
	return accountID, nil
}
