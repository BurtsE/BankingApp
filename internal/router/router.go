package router

import (
	"BankingApp/internal/config"
	"BankingApp/pkg/middleware"
	"context"
	"net/http"
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
	creditService  service.CreditService
	srv            *http.Server
}

// NewRouter — конструктор роутера
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

// InitRoutes регистрирует эндпоинты
func (r *Router) InitRoutes(userService service.UserService, bankingService service.BankingService, cardService service.CardService, creditService service.CreditService) {
	r.userService = userService
	r.bankingService = bankingService
	r.cardService = cardService
	r.creditService = creditService
	r.InitUserRoutes()
	r.InitCardRoutes()
	r.InitBankingRoutes()
	r.InitCreditRoutes()
}

// Handler возвращает основной http.Handler
func (r *Router) Handler() http.Handler {
	return r.muxRouter
}
