package router

import (
	"BankingApp/internal/config"
	"BankingApp/internal/model"
	"BankingApp/pkg/middleware"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// --------- API struct TYPES -----------

type issueCreditRequest struct {
	Amount   int64   `json:"amount"`
	Currency string  `json:"currency"`
	Months   int     `json:"months"`
	Rate     float64 `json:"rate"`
}

// ----------- HANDLERS ------------

func (r *Router) InitCreditRoutes() {
	authMiddleware := middleware.NewAuthMiddleware(config.GetJWTSecretKey())
	creditRouter := r.muxRouter.PathPrefix("/credit").Subrouter()
	creditRouter.Use(authMiddleware)

	creditRouter.HandleFunc("/issue", r.issueCreditHandler).Methods("POST")
	creditRouter.HandleFunc("/schedule", r.showScheduleHandler).Methods("POST")
	creditRouter.HandleFunc("/payment-graph", r.showPaymentsHandler).Methods("POST")
}

func (r *Router) issueCreditHandler(w http.ResponseWriter, req *http.Request) {
	userID, err := middleware.ValidateUser(req)
	if err != nil {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	var reqBody issueCreditRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if reqBody.Amount <= 0 || reqBody.Months <= 0 || reqBody.Rate <= 0 || reqBody.Currency != "RUB" {

		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	credit := model.Credit{
		UserID:      userID,
		Amount:      reqBody.Amount,
		MonthlyRate: reqBody.Rate,
		Currency:    reqBody.Currency,
		Status:      "open",
		TermMonths:  reqBody.Months,
	}
	payments, err := r.creditService.IssueCredit(ctx, credit)
	if err != nil {
		r.logger.WithError(err).Error("IssueCredit failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payments)
}
func (r *Router) showScheduleHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Автоматическое списание платежей (шедулер каждые N часов)"}`))
}

func (r *Router) showPaymentsHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Генерация графика платежей"}`))
}
