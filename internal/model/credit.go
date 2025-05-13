package model

import "time"

// Credit — кредит клиента
type Credit struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Amount     float64   `json:"amount"`
	Currency   string    `json:"currency"`
	Rate       float64   `json:"rate"`
	TermMonths int       `json:"term_months"`
	Status     string    `json:"status"` // active, closed, overdue и т.д.
	CreatedAt  time.Time `json:"created_at"`
}
