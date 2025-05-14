package app

import (
	"BankingApp/internal/config"
	"BankingApp/internal/router"
	"BankingApp/internal/service"
	"BankingApp/internal/storage"
	pgUserStorage "BankingApp/internal/storage/users/postgres"
	userService "BankingApp/internal/service/users"
	pgBankingStorage "BankingApp/internal/storage/banking/postgres"
	bankingService "BankingApp/internal/service/banking"
	"context"
	"log"

	"github.com/sirupsen/logrus"
)

type serviceProvider struct {
	cfg         *config.Config
	userStorage storage.UserStorage
	userService service.UserService
	bankingStorage storage.BankingStorage
	bankingService service.BankingService
	router      *router.Router
}

func NewSericeProvider() *serviceProvider {
	s := &serviceProvider{}
	s.Router()
	return s
}

func (s *serviceProvider) Config() *config.Config {
	if s.cfg == nil {
		cfg, err := config.InitConfig()
		if err != nil {
			log.Fatal(err)
		}
		s.cfg = cfg
	}
	return s.cfg
}

func (s *serviceProvider) UserStorage() storage.UserStorage {
	if s.userStorage == nil {
		storage, err := pgUserStorage.NewPostgresUserRepository(context.Background(), s.Config())
		if err != nil {
			log.Fatalf("could not init storage: %s", err.Error())
		}
		s.userStorage = storage
	}
	return s.userStorage
}

func (s *serviceProvider) UserService() service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(s.UserStorage(), s.Config())
	}
	return s.userService
}

func (s *serviceProvider) BankingStorage() storage.BankingStorage {
	if s.bankingStorage == nil {
		storage, err := pgBankingStorage.NewPostgresBankingRepository(context.Background(), s.Config())
		if err != nil {
			log.Fatalf("could not init storage: %s", err.Error())
		}
		s.bankingStorage = storage
	}
	return s.bankingStorage
}

func (s *serviceProvider) BankingService() service.BankingService {
	if s.bankingService == nil {
		s.bankingService = bankingService.NewBankingService(s.BankingStorage())
	}
	return s.bankingService
}


// Инициализация http-сервера.Для каждой области отдельная функция инициализации
func (s *serviceProvider) Router() *router.Router {
	if s.router == nil {
		s.router = router.NewRouter(logrus.New(), s.Config())
		s.router.InitUserRoutes(s.UserService())
		s.router.InitBankingRoutes(s.BankingService())
	}
	return s.router
}
