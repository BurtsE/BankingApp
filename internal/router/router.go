package router

import (
	"BankingApp/internal/config"
	"BankingApp/internal/router/banking"
	"BankingApp/internal/router/cards"
	"BankingApp/internal/router/user"
	"BankingApp/pkg/middleware"
	"context"
	"fmt"
	"net/http"
	"time"

	"BankingApp/internal/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Router основной роутер приложения
type Router struct {
	logger        *logrus.Logger
	muxRouter     *mux.Router
	userRouter    *user.UserSubRouter
	bankingRouter *banking.BankingSubRouter
	cardRouter    *cards.CardSubRouter
	srv           *http.Server
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
	switch {
	case r.userRouter == nil:
		return fmt.Errorf("user routes were not initialized")
	case r.bankingRouter == nil:
		return fmt.Errorf("banking routes were not initialized")
		// case r.cardRouter == nil:
		// 	return fmt.Errorf("card routes were not initialized")
	}
	return r.srv.ListenAndServe()
}

func (r *Router) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}

// InitUserService инициализирует сервис работы с пользователями
func (r *Router) InitRoutes(userService service.UserService, bankingService service.BankingService) {
	r.userRouter = user.InitUserRouter(userService, r.logger, r.muxRouter)	
	r.bankingRouter = banking.InitBankingRouter(bankingService, r.logger, r.muxRouter)
}

func (r *Router) InitCardRoutes(cardService service.CardService) {
	r.cardRouter = cards.InitCardRouter(cardService, r.logger, r.muxRouter)
}


// Handler возвращает основной http.Handler
func (r *Router) Handler() http.Handler {
	return r.muxRouter
}
