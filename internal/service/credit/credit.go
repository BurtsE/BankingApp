package credit

import (
	"BankingApp/internal/model"
	"BankingApp/internal/storage"
	"context"
	"fmt"
	"math"
	"time"
)

type CreditService struct {
	storage storage.CreditStorage
}

func NewCreditService(storage storage.CreditStorage) *CreditService {
	service := &CreditService{storage: storage}
	return service
}

func (cs *CreditService) IssueCredit(ctx context.Context, credit model.Credit) ([]*model.PaymentSchedule, error) {
	tx, err := cs.storage.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	payments := make([]*model.PaymentSchedule, credit.TermMonths)
	creditID, err := cs.storage.IssueCredit(ctx, credit)
	if err != nil {
		return nil, err
	}
	amount := calculatePayment(float64(credit.Amount), credit.MonthlyRate, credit.TermMonths)
	currTime := time.Now()
	for i := range credit.TermMonths {
		payments[i] = &model.PaymentSchedule{
			CreditID: creditID,
			DueDate:  currTime.AddDate(0, i+1, 0),
			Amount:   amount,
			IsPaid:   false,
			PaidAt:   nil,
		}
		payments[i].ID, err = cs.storage.IssuePayment(ctx, *payments[i])
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit error: %w", err)
	}
	return payments, nil
}

func calculatePayment(sum, rate float64, months int) float64 {
	p := math.Pow(rate+1, float64(months))
	return sum * rate * p / (p - 1)
}
