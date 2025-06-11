package postgres

import (
	"BankingApp/internal/model"
	"context"
	"fmt"
)

func (p *PostgresRepository) IssueCredit(ctx context.Context, credit model.Credit) (int64, error) {
	const query = `
		INSERT 
			INTO credits (user_id, amount, currency, monthly_rate, term_months, status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var creditID int64
	err := p.pool.QueryRow(ctx, query,
		credit.UserID,
		credit.Amount,
		credit.Currency,
		credit.MonthlyRate,
		credit.TermMonths,
		credit.Status,
		credit.CreatedAt,
	).Scan(&creditID)

	if err != nil {
		return 0, fmt.Errorf("failed to issue credit in transaction: %w", err)
	}

	return creditID, nil
}

func (p *PostgresRepository) IssuePayment(ctx context.Context, payment model.PaymentSchedule) (int64, error) {
	const query = `
		INSERT 
			INTO payment_schedules (credit_id, payment_num, due_date, amount, is_paid, paid_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var paymentID int64
	err := p.pool.QueryRow(ctx, query,
		payment.CreditID,
		payment.PaymentNum,
		payment.DueDate,
		payment.Amount,
		payment.IsPaid,
		payment.PaidAt,
	).Scan(&paymentID)

	if err != nil {
		return 0, fmt.Errorf("failed to issue payment: %w", err)
	}

	return paymentID, nil
}
