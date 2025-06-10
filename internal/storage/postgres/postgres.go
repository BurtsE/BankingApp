package postgres

import (
	"BankingApp/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	encryptionPublicKey  []byte
	encryptionPrivateKey []byte
	pool                 *pgxpool.Pool
}

// NewPostgresUserRepository конструктор для pgxpool
func NewPostgresRepository(ctx context.Context, cfg *config.Config) (*PostgresRepository, error) {
	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database)
	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	repo := &PostgresRepository{
		pool:                 pool,
		encryptionPublicKey:  []byte(config.GetEncryptionPublicKey()),
		encryptionPrivateKey: []byte(config.GetEncryptionPrivateKey()),
	}

	return repo, nil
}

func (p *PostgresRepository) Close() {
	p.pool.Close()
}
