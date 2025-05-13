
package model

// Transaction — история операций
type Transaction struct {
	ID              int64     `json:"id"`
	AccountID       int64     `json:"account_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Type            string    `json:"type"`             // deposit, withdraw, transfer, payment и т.д.
	Status          string    `json:"status"`           // success, pending, failed
	Description     string    `json:"description"`
	RelatedEntityID *int64    `json:"related_entity_id,omitempty"` // напр: ID карты, кредита, если нужно
	CreatedAt       time.Time `json:"created_at"`
}