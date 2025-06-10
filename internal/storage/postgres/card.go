package postgres

import (
	"BankingApp/internal/model"
	"context"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func (p *PostgresRepository) CreateVirtualCard(ctx context.Context, card *model.Card) (int64, error) {
	encryptedPan, err := p.encryptWithPGP([]byte(card.PAN))
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	query := `
		INSERT INTO cards (account_id, encrypted_pan, expiry_month, expiry_year , cardholder_name, is_active, created_at)
		VALUES ($1, $2::bytea, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id int64
	err = p.pool.QueryRow(ctx, query, &card.AccountID, &encryptedPan, &card.ExpiryMonth, &card.ExpiryYear, &card.CardholderName, &card.IsActive, &card.CreatedAt).Scan(&id)
	return id, err

}

func (p *PostgresRepository) GetCardsByAccount(ctx context.Context, accountID int64) ([]*model.Card, error) {
	query := `
		SELECT id, encode(encrypted_pan, 'escape')::text, expiry_month, expiry_year , cardholder_name, is_active, created_at
		FROM cards 
		WHERE account_id=$1
	`
	rows, err := p.pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("GetCardsByAccount: %w", err)
	}
	defer rows.Close()

	cards := make([]*model.Card, 0)
	for rows.Next() {
		var (
			card          model.Card
			encrypted_pan []byte
		)
		if err := rows.Scan(&card.ID, &encrypted_pan, &card.ExpiryMonth, &card.ExpiryYear, &card.CardholderName, &card.IsActive, &card.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCardsByAccount scan: %w", err)
		}
		pan, err := p.decryptWithPGP(encrypted_pan)
		if err != nil {
			return nil, fmt.Errorf("GetCardsByAccount decryption: %w", err)
		}
		card.PAN = string(pan)
		card.AccountID = accountID
		cards = append(cards, &card)
	}
	return cards, rows.Err()
}

// Encrypt plaintext message using a public key
func (p *PostgresRepository) encryptWithPGP(data []byte) (string, error) {
	pgp := crypto.PGP()
	// Encrypt plaintext message using a public key
	encHandle, err := pgp.Encryption().Recipient(p.publicKey).New()
	if err != nil {
		return "", err
	}
	pgpMessage, err := encHandle.Encrypt(data)
	if err != nil {
		return "", err
	}
	armored, err := pgpMessage.Armor()
	if err != nil {
		return "", err
	}
	return armored, nil
}

// Decrypt armored encrypted message using the private key and obtain the plaintext
func (p *PostgresRepository) decryptWithPGP(encrypted []byte) (string, error) {
	pgp := crypto.PGP()

	// Decrypt armored encrypted message using the private key and obtain the plaintext
	decHandle, err := pgp.Decryption().DecryptionKey(p.privateKey).New()
	if err != nil {
		return "", err
	}
	decrypted, err := decHandle.Decrypt(encrypted, crypto.Armor)
	if err != nil {
		return "", err
	}
	return decrypted.String(), nil
}
