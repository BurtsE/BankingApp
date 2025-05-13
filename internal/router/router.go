
package router

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Router основной роутер приложения
type Router struct {
	muxRouter *mux.Router
}

// NewRouter — конструктор роутера, регистрирует все endpoint'ы
func NewRouter() *Router {
	r := &Router{
		muxRouter: mux.NewRouter(),
	}
	r.initPublicRoutes()
	r.initProtectedRoutes()
	return r
}

// Handler возвращает основной http.Handler
func (r *Router) Handler() http.Handler {
	return r.muxRouter
}

// --- PUBLIC ROUTES ---
func (r *Router) initPublicRoutes() {
	r.muxRouter.HandleFunc("/register", registerHandler).Methods("POST")
	r.muxRouter.HandleFunc("/login", loginHandler).Methods("POST")
}

// --- PROTECTED ROUTES (JWT Auth Required) ---
func (r *Router) initProtectedRoutes() {
	secured := r.muxRouter.NewRoute().Subrouter()
	secured.Use(jwtMiddleware) // JWT middleware

	secured.HandleFunc("/accounts", createAccountHandler).Methods("POST")
	secured.HandleFunc("/cards", issueCardHandler).Methods("POST")
	secured.HandleFunc("/transfer", transferHandler).Methods("POST")
	secured.HandleFunc("/analytics", analyticsHandler).Methods("GET")
	secured.HandleFunc("/credits/{creditId}/schedule", creditScheduleHandler).Methods("GET")
	secured.HandleFunc("/accounts/{accountId}/predict", predictBalanceHandler).Methods("GET")
}

// --- JWT Middleware (заглушка) ---
func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		// Здесь будет полноценная JWT валидация
		if tokenStr == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// TODO: validate JWT (пока всегда true)
		next.ServeHTTP(w, r)
	})
}

// --- HANDLERS (заглушки для примера) ---

// Публичные
func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Пользователь успешно зарегистрирован"}`))
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"fake-jwt-token"}`))
}

// Защищённые
func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Счёт создан"}`))
}
func issueCardHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Карта выпущена"}`))
}
func transferHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Перевод выполнен"}`))
}
func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"analytics":"stub data"}`))
}
func creditScheduleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"schedule":"stub data"}`))
}
func predictBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"prediction":"stub data"}`))
}
