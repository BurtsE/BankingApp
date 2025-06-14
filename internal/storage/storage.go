package storage

import (
	"BankingApp/internal/model"
	"context"
)

// UserRepository — интерфейс взаимодействия с таблицей пользователей
type UserStorage interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, userID string) (*model.User, error)
}

type BankingStorage interface {
	BeginTransaction(ctx context.Context) (Transaction, error)
	CreateAccount(ctx context.Context, userID string, currency string) (*model.Account, error)
	GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error)
	GetAccountsByUser(ctx context.Context, userID string) ([]*model.Account, error)
	UpdateAccountBalance(ctx context.Context, accountID int64, amount float64) error
}

type CardStorage interface {
	CreateVirtualCard(ctx context.Context, card *model.Card) (int64, error)
	GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error)
	GetCardsByAccount(ctx context.Context, accountID int64) ([]*model.Card, error)
}

type CreditStorage interface {
	BeginTransaction(ctx context.Context) (Transaction, error)
	IssueCredit(ctx context.Context, credit model.Credit) (int64, error)
	IssuePayment(ctx context.Context, payment model.PaymentSchedule) (int64, error)
}

type Transaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}
