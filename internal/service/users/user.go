package users

import (
	"BankingApp/internal/config"
	"BankingApp/internal/model"
	"BankingApp/internal/service"
	"BankingApp/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var _ service.UserService = (*Service)(nil)

type Service struct {
	repo      storage.UserStorage
	jwtSecret []byte
}

func NewUserService(repo storage.UserStorage, cfg *config.Config) *Service {
	return &Service{repo: repo, jwtSecret: []byte(config.GetJWTSecretKey())}
}

func (s *Service) Register(ctx context.Context, email, username, password, fullName string) (*model.User, error) {
	// Проверим уникальность email и username
	if user, err := s.repo.FindByEmail(ctx, email); user != nil {
		return nil, errors.New("email уже зарегистрирован")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	if user, err := s.repo.FindByUsername(ctx, username); user != nil {
		return nil, errors.New("username уже занят")
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		UUID:      uuid.New().String(),
		Email:     email,
		Password:  string(hash),
		FullName:  fullName,
		CreatedAt: time.Now(),
	}
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Authenticate(ctx context.Context, email, password string) (string, time.Time, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", time.Time{}, errors.New("неверный логин или пароль")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", time.Time{}, errors.New("неверный логин или пароль")
	}

	// Генерируем JWT
	exp := time.Now().Add(24 * time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   user.UUID,
		ExpiresAt: jwt.NewNumericDate(exp),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenStr, exp, nil
}

func (s *Service) GetByID(ctx context.Context, userID int64) (*model.User, error) {
	return s.repo.FindByID(ctx, userID)
}
