package postgres

import (
	"BankingApp/internal/model"
	"bytes"
	"context"
	"fmt"
	"io"

	"golang.org/x/crypto/openpgp"
)

func (p *PostgresRepository) CreateVirtualCard(ctx context.Context, card *model.Card) (int64, error) {
	encryptedPan, err := p.encryptWithPGP([]byte(card.PAN))
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	query := `
		INSERT INTO cards (account_id, encrypted_pan, expiry_month, expiry_year , cardholder_name, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id int64
	err = p.pool.QueryRow(ctx, query, &card.AccountID, &encryptedPan, &card.ExpiryMonth, &card.ExpiryYear, &card.CardholderName, &card.IsActive, &card.CreatedAt).Scan(&id)
	return id, err

}

func (p *PostgresRepository) GetCardsByAccount(ctx context.Context, accountID int64) ([]*model.Card, error) {
	query := `
		SELECT id, encrypted_pan, expiry_month, expiry_year , cardholder_name, is_active, created_at
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
		if err := rows.Scan(&card.ID, &encrypted_pan, &card.ExpiryMonth, card.ExpiryYear, &card.CardholderName, &card.IsActive, &card.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCardsByAccount scan: %w", err)
		}
		pan, err := p.decryptWithPGP(encrypted_pan)
		if err != nil {
			return nil, fmt.Errorf("GetCardsByAccount decryption: %w", err)
		}
		card.PAN = string(pan)
		cards = append(cards, &card)
	}
	return cards, rows.Err()
}

func (p *PostgresRepository) encryptWithPGP(data []byte) ([]byte, error) {
	// Загружаем публичный ключ
	keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(p.encryptionPrivateKey))
	if err != nil {
		return nil, err
	}

	// Шифруем
	var buf bytes.Buffer
	writer, err := openpgp.Encrypt(&buf, keyring, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	writer.Close()

	return buf.Bytes(), nil
}

func (p *PostgresRepository) decryptWithPGP(encrypted []byte) ([]byte, error) {
	// Загружаем приватный ключ
	keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(p.encryptionPrivateKey))
	if err != nil {
		return nil, err
	}

	// Расшифровываем
	md, err := openpgp.ReadMessage(bytes.NewReader(encrypted), keyring, nil, nil)
	if err != nil {
		return nil, err
	}

	decrypted, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
