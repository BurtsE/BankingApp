package postgres

import (
	"BankingApp/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	query := `
		INSERT INTO users (email, password, full_name, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var id int64
	err := r.pool.QueryRow(ctx, query, user.Email, user.Password, user.FullName, user.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password, full_name, created_at
		FROM users
		WHERE email = $1
		LIMIT 1;
	`
	row := r.pool.QueryRow(ctx, query, email)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.CreatedAt)
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
		SELECT id, email, password, full_name, created_at
		FROM users
		WHERE full_name = $1
		LIMIT 1;
	`
	row := r.pool.QueryRow(ctx, query, username)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.CreatedAt)
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
		SELECT id, email, password, full_name, created_at
		FROM users
		WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, userID)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
