package postgres

import (
	"BankingApp/internal/config"
	"BankingApp/internal/model"
	"BankingApp/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresBankingRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresBankingRepository(ctx context.Context, cfg *config.Config) (*PostgresBankingRepository, error) {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.BankingPostgres.Host, cfg.BankingPostgres.Username, cfg.BankingPostgres.Password, cfg.BankingPostgres.Database)
	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &PostgresBankingRepository{pool: pool}, nil
}

func (p *PostgresBankingRepository) CreateAccount(ctx context.Context, userID int64, currency string) (*model.Account, error) {
	query := "INSERT INTO accounts (user_id, currency, balance) VALUES ($1, $2, 0.0) RETURNING id, user_id, currency, balance"
	var acc model.Account
	err := p.pool.QueryRow(ctx, query, userID, currency).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Currency,
		&acc.Balance,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateAccount: %w", err)
	}
	return &acc, nil
}

func (p *PostgresBankingRepository) GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error) {
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

func (p *PostgresBankingRepository) GetAccountsByUser(ctx context.Context, userID int64) ([]*model.Account, error) {
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

func (p *PostgresBankingRepository) UpdateAccountBalance(ctx context.Context, accountID int64, amount float64) error {
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

func (p *PostgresBankingRepository) BeginTransaction(ctx context.Context) (storage.Transaction, error) {
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("BeginTransaction: %w", err)
	}
	return tx, nil
}

