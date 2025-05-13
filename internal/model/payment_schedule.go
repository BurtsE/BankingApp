package model

import "time"

// PaymentSchedule — график платежей по кредиту
type PaymentSchedule struct {
	ID         int64      `json:"id"`
	CreditID   int64      `json:"credit_id"`
	PaymentNum int        `json:"payment_num"`
	DueDate    time.Time  `json:"due_date"`
	Amount     float64    `json:"amount"`
	IsPaid     bool       `json:"is_paid"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
}
