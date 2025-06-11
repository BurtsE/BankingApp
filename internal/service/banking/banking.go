package banking

import (
	"BankingApp/internal/model"
	"BankingApp/internal/storage"
	"context"
	"errors"
	"fmt"
)

type BankingService struct {
	storage storage.BankingStorage
}

func NewBankingService(storage storage.BankingStorage) *BankingService {
	return &BankingService{storage: storage}
}

func (s *BankingService) CreateAccount(ctx context.Context, userID string, currency string) (*model.Account, error) {
	return s.storage.CreateAccount(ctx, userID, currency)
}

func (s *BankingService) Deposit(ctx context.Context, accountID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return s.storage.UpdateAccountBalance(ctx, accountID, amount)
}

func (s *BankingService) Withdraw(ctx context.Context, accountID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	account, err := s.storage.GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}
	if account.Balance < amount {
		return errors.New("insufficient balance")
	}
	return s.storage.UpdateAccountBalance(ctx, accountID, -amount)
}

func (s *BankingService) Transfer(ctx context.Context, fromAccountID, toAccountID int64, amount float64) error {
	if fromAccountID == toAccountID {
		return errors.New("cannot transfer to the same account")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	from, err := s.storage.GetAccountByID(ctx, fromAccountID)
	if err != nil {
		return err
	}
	_, err = s.storage.GetAccountByID(ctx, toAccountID)
	if err != nil {
		return err
	}
	if from.Balance < amount {
		return errors.New("insufficient balance")
	}
	tx, err := s.storage.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := s.storage.UpdateAccountBalance(ctx, fromAccountID, -amount); err != nil {
		return err
	}
	if err := s.storage.UpdateAccountBalance(ctx, toAccountID, amount); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}

func (s *BankingService) GetAccountsByUser(ctx context.Context, userID string) ([]*model.Account, error) {
	return s.storage.GetAccountsByUser(ctx, userID)
}

func (s *BankingService) GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error) {
	return s.storage.GetAccountByID(ctx, accountID)
}
