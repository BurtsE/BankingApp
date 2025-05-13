package model


// Card — карта, привязанная к счету (CVV зашифрован)
type Card struct {
	ID             int64     `json:"id"`
	AccountID      int64     `json:"account_id"`
	Number         string    `json:"number"`
	ExpiryMonth    int       `json:"expiry_month"`
	ExpiryYear     int       `json:"expiry_year"`
	EncryptedCVV   string    `json:"encrypted_cvv"`     // зашифрованное значение!
	CardholderName string    `json:"cardholder_name"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}