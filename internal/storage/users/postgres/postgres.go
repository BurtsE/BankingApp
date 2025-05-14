package repository

import (
	"BankingApp/internal/config"
	"BankingApp/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository конструктор для pgxpool
func NewPostgresUserRepository(ctx context.Context, cfg *config.Config) (*PostgresUserRepository, error) {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.UserPostgres.Host, cfg.UserPostgres.Username, cfg.UserPostgres.Password, cfg.UserPostgres.Database)
	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &PostgresUserRepository{
		pool: pool,
	}, nil
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *model.User) (int64, error) {
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

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
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

func (r *PostgresUserRepository) FindByID(ctx context.Context, userID int64) (*model.User, error) {
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
