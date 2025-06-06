package model

import "time"

// Account — банковский счет пользователя
type Account struct {
	ID        int64     `json:"id"`
	Balance   float64   `json:"balance"`
	UserID    string    `json:"user_id"`
	Number    string    `json:"number"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}
