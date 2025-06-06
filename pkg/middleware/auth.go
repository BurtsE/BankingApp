package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
		claims.GetIssuer()
		UUID, err := claims.GetSubject()
		if err != nil {
			return "", errors.New("user_id missing in token")
		}
		return UUID, nil
	}
	return "", errors.New("invalid token")
}
