package router

import (
	"BankingApp/internal/config"
	"BankingApp/pkg/middleware"
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"BankingApp/internal/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Router основной роутер приложения
type Router struct {
	logger         *logrus.Logger
	muxRouter      *mux.Router
	userService    service.UserService
	bankingService service.BankingService
	cardService    service.CardService
	srv            *http.Server
}

// NewRouter — конструктор роутера, регистрирует все endpoint'ы
func NewRouter(logger *logrus.Logger, cfg *config.Config) *Router {
	r := &Router{
		muxRouter: mux.NewRouter().PathPrefix("/api/v1").Subrouter(),
		logger:    logger,
	}
	r.srv = &http.Server{
		Handler:      r.Handler(),
		Addr:         ":" + cfg.ServerPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	r.muxRouter.Use(loggingMiddleware)
	return r
}
func (r *Router) Start() error {
	return r.srv.ListenAndServe()
}

func (r *Router) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}

// InitUserService инициализирует сервис работы с пользователями
func (r *Router) InitRoutes(userService service.UserService, bankingService service.BankingService, cardService service.CardService) {
	r.userService = userService
	r.bankingService = bankingService
	r.cardService = cardService
	r.InitUserRoutes()
	r.InitCardRoutes()
	r.InitBankingRoutes()
}

// Handler возвращает основной http.Handler
func (r *Router) Handler() http.Handler {
	return r.muxRouter
}

func ParseIDFromVars(r *http.Request, varName string) (int64, error) {
	vars := mux.Vars(r)
	raw, ok := vars[varName]
	if !ok {
		return 0, errors.New("missing id")
	}
	return strconv.ParseInt(raw, 10, 64)
}
