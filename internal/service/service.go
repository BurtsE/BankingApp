package service

import (
	"context"
	"time"

	"BankingApp/internal/model"
)

type UserService interface {
	Register(ctx context.Context, email, username, password, fullName string) (*model.User, error)
	Authenticate(ctx context.Context, email, password string) (jwtToken string, expiresAt time.Time, err error)
	GetByID(ctx context.Context, userID string) (*model.User, error)
}

type BankingService interface {
	CreateAccount(ctx context.Context, userID string, currency string) (*model.Account, error)
	Deposit(ctx context.Context, accountID int64, amount float64) error
	Withdraw(ctx context.Context, accountID int64, amount float64) error
	Transfer(ctx context.Context, fromAccountID, toAccountID int64, amount float64) error
	GetAccountsByUser(ctx context.Context, userID string) ([]*model.Account, error)
	GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error)
}

type CardService interface {
	GenerateVirtualCard(ctx context.Context, accountID int64, cardholderName string) (*model.Card, error)
	GetCardsByAccount(ctx context.Context, accountID int64) ([]*model.Card, error)
	GetCardByIDForOwner(ctx context.Context, cardID, ownerUserID int64) (*model.Card, error) // с расшифровкой
}

type CreditService interface {
	IssueCredit(ctx context.Context, credit model.Credit) ([]*model.PaymentSchedule, error)
	// GeneratePaymentSchedule(ctx context.Context, creditID int64) ([]*model.PaymentSchedule, error)
	// AutoWithdrawCreditPayments(ctx context.Context, now time.Time) error // для background-шейдулера
	// ProcessFine(ctx context.Context, creditID int64, overdueAmount float64) error
	// GetCreditsByUser(ctx context.Context, userID int64) ([]*model.Credit, error)
	// GetPaymentSchedule(ctx context.Context, creditID int64) ([]*model.PaymentSchedule, error)
}

type AnalyticService interface {
	// GetMonthlyReport(ctx context.Context, userID int64, month time.Month, year int) (*model.Analytics, error)
	// GetCreditLoadAnalytics(ctx context.Context, userID int64) (*model.CreditLoadAnalytics, error)
	// PredictBalance(ctx context.Context, accountID int64, days int) (*model.Prediction, error)
}

type IntegrationService interface {
	GetCentralBankKeyRate(ctx context.Context, date time.Time) (float64, error)       // SOAP запрос и парсинг XML
	SendPaymentNotificationEmail(ctx context.Context, to, subject, body string) error // SMTP/SIMPLE
}
