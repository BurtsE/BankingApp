package cards

import (
	"BankingApp/internal/config"
	"BankingApp/internal/model"
	"BankingApp/internal/storage"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"log"
	"math/big"
	"time"

	"crypto/rand"
)

type CardService struct {
	storage storage.CardStorage
	gcm     cipher.AEAD
}

func NewCardService(storage storage.CardStorage) *CardService {
	service := &CardService{storage: storage}
	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher([]byte(config.GetCardSecretKey()))
	if err != nil {
		log.Fatalf("%v", err)
	}
	// gcm or Galois/Counter Mode, is a mode of operation
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("%v", err)
	}
	service.gcm = gcm
	return service
}

func (cs *CardService) GenerateVirtualCard(ctx context.Context, accountID int64, cardholderName string) (*model.Card, error) {
	pan, err := cs.generatePAN()
	if err != nil {
		return nil, err
	}
	card := &model.Card{
		AccountID:      accountID,
		CardholderName: cardholderName,
		PAN:            pan,
		IsActive:       true,
		CreatedAt:      time.Now(),
	}
	card.GenerateTimeExpiry()
	card.CVV = cs.generateCVV(card.PAN, time.Date(card.ExpiryYear, time.Month(card.ExpiryMonth), 0, 0, 0, 0, 0, time.Local))
	id, err := cs.storage.CreateVirtualCard(ctx, card)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании карты: %w", err)
	}
	card.ID = id
	return card, nil
}
func (cs *CardService) GetCardsByAccount(ctx context.Context, accountID int64) ([]*model.Card, error) {
	cards, err := cs.storage.GetCardsByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения карт: %w", err)
	}
	for _, card := range cards {
		card.CVV = cs.generateCVV(card.PAN, time.Date(card.ExpiryYear, time.Month(card.ExpiryMonth), 0, 0, 0, 0, 0, time.Local))
	}
	return cards, nil
}
func (cs *CardService) GetCardByIDForOwner(ctx context.Context, cardID, ownerUserID int64) (*model.Card, error) {
	panic("implement")
}

func (cs *CardService) generateCVV(pan string, expiryDate time.Time) string {
	dst := []byte{}
	nonce := make([]byte, cs.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	seal := cs.gcm.Seal(dst, nonce, []byte(pan+expiryDate.String()), nil)[:3]
	cvv := make([]byte, 0, 3)
	for i := range 3 {
		cvv = append(cvv, seal[i]%10)
	}
	return string(cvv)
}

func (cs *CardService) generatePAN() (string, error) {
	pan := make([]byte, 15)
	pan[0] = byte(4)
	sum := 8
	for i := range 14 {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("ошибка генерации карты: %v", err)
		}
		pan[i+1] = byte(digit.Int64())
		if i%2 != 0 {
			sum += int(digit.Int64()) * 2
		} else {
			sum += int(digit.Int64())
		}
	}
	lastDigit := 0
	for sum+lastDigit%10 != 0 {
		lastDigit++
	}
	pan = append(pan, byte(lastDigit))
	return string(pan), nil
}
