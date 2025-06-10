package model

import (
	"time"
)

// Card — карта, привязанная к счету (CVV зашифрован)
type Card struct {
	ExpiryMonth    int       `json:"expiry_month"`
	ExpiryYear     int       `json:"expiry_year"`
	ID             int64     `json:"id"`
	AccountID      int64     `json:"account_id"`
	PAN            string    `json:"number"` // При сохранении шифруем
	CVV            string    `json:"cvv"`    // Не храним
	CardholderName string    `json:"cardholder_name"`
	CreatedAt      time.Time `json:"created_at"`
	IsActive       bool      `json:"is_active"`
}

func (c *Card) GenerateTimeExpiry() {
	c.ExpiryYear = time.Now().AddDate(4, 0, 0).Year()
	c.ExpiryMonth = int(time.Now().Month())
}
