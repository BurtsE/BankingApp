package model

import "time"

// User — данные пользователя
type User struct {
	UUID      string    `json:"uuid"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // bcrypt hash
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}
