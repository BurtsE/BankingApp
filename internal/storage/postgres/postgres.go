package postgres

import (
	"BankingApp/internal/config"
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/jackc/pgx/v5/pgxpool"
)

const passPhrase = "pass_phrase"

type PostgresRepository struct {
	config     *packet.Config
	privateKey *crypto.Key
	publicKey  *crypto.Key
	pool       *pgxpool.Pool
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
	// Configure key settings
	pgpConfig := &packet.Config{
		Algorithm:     packet.PubKeyAlgoRSA,
		RSABits:       4096,
		DefaultCipher: packet.CipherAES256,
		Rand:          rand.Reader,
		Time:          func() time.Time { return time.Now().UTC() },
	}

	publicKey, err := crypto.NewKeyFromArmored(config.GetEncryptionPublicKey())
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.NewPrivateKeyFromArmored(config.GetEncryptionPrivateKey(), []byte(passPhrase))
	if err != nil {
		return nil, err
	}
	repo := &PostgresRepository{
		pool:       pool,
		config:     pgpConfig,
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	return repo, nil
}

func (p *PostgresRepository) Close() {
	p.pool.Close()
}
