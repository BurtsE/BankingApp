package postgres

import (
	"BankingApp/internal/model"
	"BankingApp/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (p *PostgresRepository) CreateAccount(ctx context.Context, userID string, currency string) (*model.Account, error) {
	query := "INSERT INTO accounts (user_id, currency, balance) VALUES ($1, $2, 0.0) RETURNING id, user_id, currency, balance, is_active"
	var acc model.Account
	err := p.pool.QueryRow(ctx, query, userID, currency).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Currency,
		&acc.Balance,
		&acc.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateAccount: %w", err)
	}
	return &acc, nil
}

func (p *PostgresRepository) GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error) {
	query := "SELECT id, user_id, currency, balance FROM accounts WHERE id=$1"
	var acc model.Account
	err := p.pool.QueryRow(ctx, query, accountID).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Currency,
		&acc.Balance,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("GetAccountByID: %w", err)
	}
	return &acc, nil
}

func (p *PostgresRepository) GetAccountsByUser(ctx context.Context, userID string) ([]*model.Account, error) {
	query := "SELECT id, user_id, currency, balance FROM accounts WHERE user_id=$1"
	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetAccountsByUser: %w", err)
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.Currency, &acc.Balance); err != nil {
			return nil, fmt.Errorf("GetAccountsByUser scan: %w", err)
		}
		accounts = append(accounts, &acc)
	}
	return accounts, rows.Err()
}

func (p *PostgresRepository) UpdateAccountBalance(ctx context.Context, accountID int64, amount float64) error {
	// This is a simple example. In practice you want to check for negative balances in a transaction!
	query := "UPDATE accounts SET balance = balance + $1 WHERE id = $2"
	result, err := p.pool.Exec(ctx, query, amount, accountID)
	if err != nil {
		return fmt.Errorf("UpdateAccountBalance: %w", err)
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return errors.New("account not found")
	}
	return nil
}

func (p *PostgresRepository) BeginTransaction(ctx context.Context) (storage.Transaction, error) {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("BeginTransaction: %w", err)
	}
	return tx, nil
}
