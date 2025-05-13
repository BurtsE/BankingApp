package repository

import (
	"BankingApp/internal/model"
	"context"
)

// UserRepository — интерфейс взаимодействия с таблицей пользователей
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, userID int64) (*model.User, error)
}
