package app

import (
	"BankingApp/internal/config"
	"BankingApp/internal/router"
	"BankingApp/internal/service"
	bankingService "BankingApp/internal/service/banking"
	cardService "BankingApp/internal/service/cards"
	userService "BankingApp/internal/service/users"
	storageImpl "BankingApp/internal/storage/postgres"
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// serviceProvider есть DI-контейнер, производящий инициализацию сервисов и подключения к БД.
// Глобальный контекст добавлен для Graceful shutdown
type serviceProvider struct {
	cfg            *config.Config
	storage        *storageImpl.PostgresRepository
	userService    service.UserService
	bankingService service.BankingService
	cardService    service.CardService
	router         *router.Router
	logger         *logrus.Logger
	errG           *errgroup.Group
	ctx            context.Context
}

func NewServiceProvider(ctx context.Context) *serviceProvider {
	errG, gCtx := errgroup.WithContext(ctx)
	s := &serviceProvider{
		errG: errG,
		ctx:  gCtx,
	}

	s.Router()
	return s
}

func (s *serviceProvider) Start() {
	s.errG.Go(func() error {
		s.logger.Printf("starting server on port: %s", s.Config().ServerPort)
		return s.router.Start()
	})
	if err := s.errG.Wait(); err != nil {
		s.logger.Printf("exit reason: %s \n", err)
	}
}

func (s *serviceProvider) Config() *config.Config {
	if s.cfg == nil {
		cfg, err := config.InitConfig()
		if err != nil {
			s.logger.Fatal(err)
		}
		s.cfg = cfg
	}
	return s.cfg
}
func (s *serviceProvider) Storage() *storageImpl.PostgresRepository {
	if s.storage == nil {
		storage, err := storageImpl.NewPostgresRepository(context.Background(), s.Config())
		if err != nil {
			s.logger.Fatalf("could not init storage: %s", err.Error())
		}
		s.storage = storage
		s.errG.Go(func() error {
			<-s.ctx.Done()
			s.logger.Println("closing user database...")
			s.storage.Close()
			return nil
		})
	}
	return s.storage
}

func (s *serviceProvider) UserService() service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(s.Storage(), s.Config())
	}
	return s.userService
}

func (s *serviceProvider) BankingService() service.BankingService {
	if s.bankingService == nil {
		s.bankingService = bankingService.NewBankingService(s.Storage())
	}
	return s.bankingService
}

func (s *serviceProvider) CardService() service.CardService {
	if s.cardService == nil {
		s.cardService = cardService.NewCardService(s.Storage())
	}
	return s.cardService
}

// Инициализация http-сервера.Для каждой области отдельная функция инициализации
func (s *serviceProvider) Router() *router.Router {
	if s.router == nil {
		s.router = router.NewRouter(s.Logger(), s.Config())
		s.router.InitRoutes(s.UserService(), s.BankingService(), s.CardService())
		s.errG.Go(func() error {
			<-s.ctx.Done()
			s.logger.Println("shutting down server...")
			return s.router.Stop(s.ctx)
		})
	}
	return s.router
}
func (s *serviceProvider) Logger() *logrus.Logger {
	if s.logger == nil {
		s.logger = logrus.New()
		if s.Config().LogLevel == "DEBUG" {
			s.logger.Level = logrus.DebugLevel
		}
	}
	return s.logger
}
