package cards

import (
	"BankingApp/internal/config"
	"BankingApp/internal/model"
	"BankingApp/internal/storage"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
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
	cvv, err := cs.generateCVV(card)
	if err != nil {
		return nil, err
	}
	card.CVV = cvv

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
		cvv, err := cs.generateCVV(card)
		if err != nil {
			return nil, err
		}
		card.CVV = cvv
	}
	return cards, nil
}
func (cs *CardService) GetCardByIDForOwner(ctx context.Context, cardID, ownerUserID int64) (*model.Card, error) {
	panic("implement")
}

func (cs *CardService) generateCVV(card *model.Card) (string, error) {
	// Create a hash of the combined input
	combined := card.PAN + "|" + time.Date(card.ExpiryYear, time.Month(card.ExpiryMonth), 0, 0, 0, 0, 0, time.Local).String()

	// Hash the combined input to get a consistent key
	hasher := sha256.New()
	hasher.Write([]byte(combined))
	key := hasher.Sum(nil)[:aes.BlockSize] // Use first 16 bytes as AES key

	// Create a cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// We'll use a fixed nonce (since input is always the same)
	nonce := make([]byte, gcm.NonceSize())

	// Encrypt some fixed data (we just need consistent output)
	// Using the variables again as plaintext to incorporate all input
	ciphertext := gcm.Seal(nil, nonce, []byte(combined), nil)

	// Convert first 4 bytes of ciphertext to a number
	num := binary.BigEndian.Uint32(ciphertext[:4])

	// Reduce to three digits (000-999)

	cvv := make([]byte, 0, 3)
	for range 3 {
		cvv = append(cvv, byte(num%10)+48)
		num /= 10
	}
	return string(cvv), nil
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
	for (sum+lastDigit)%10 != 0 {
		lastDigit++
	}
	pan = append(pan, byte(lastDigit))
	for i := range pan {
		pan[i] += 48
	}
	return string(pan), nil
}
