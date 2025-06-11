package model

import "time"

// Credit — кредит клиента
type Credit struct {
	ID          int64     `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      int64     `json:"amount"`
	Currency    string    `json:"currency"`
	MonthlyRate float64   `json:"rate"`
	TermMonths  int       `json:"term_months"`
	Status      string    `json:"status"` // active, closed, overdue и т.д.
	CreatedAt   time.Time `json:"created_at"`
}
