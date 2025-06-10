package postgres

import (
	"BankingApp/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (uuid, email, password, full_name, created_at)
		VALUES ($1, $2, $3, $4, $5)

	`
	_, err := r.pool.Exec(ctx, query, user.UUID, user.Email, user.PasswordHash, user.FullName, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT uuid, email, password, full_name, created_at
		FROM users
		WHERE email = $1
		LIMIT 1;
	`
	row := r.pool.QueryRow(ctx, query, email)
	user := &model.User{}
	err := row.Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.FullName, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT uuid, email, password, full_name, created_at
		FROM users
		WHERE full_name = $1
		LIMIT 1;
	`
	row := r.pool.QueryRow(ctx, query, username)
	user := &model.User{}
	err := row.Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.FullName, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, userID int64) (*model.User, error) {
	query := `
		SELECT uuid, email, password, full_name, created_at
		FROM users
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, userID)
	user := &model.User{}
	err := row.Scan(&user.UUID, &user.Email, &user.PasswordHash, &user.FullName, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
