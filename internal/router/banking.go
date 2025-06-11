package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"BankingApp/internal/config"
	"BankingApp/pkg/middleware"

	"github.com/gorilla/mux"
)

// --- PROTECTED ROUTES (JWT Auth Required) ---

func (r *Router) InitBankingRoutes() {
	authMiddleware := middleware.NewAuthMiddleware(config.GetJWTSecretKey())
	bankingRouter := r.muxRouter.PathPrefix("/banking").Subrouter()
	bankingRouter.Use(authMiddleware)
	bankingRouter.HandleFunc("/account", r.createAccountHandler).Methods("POST")
	bankingRouter.HandleFunc("/account/{id:[0-9]+}/deposit", r.depositHandler).Methods("POST")
	bankingRouter.HandleFunc("/account/{id:[0-9]+}/withdraw", r.withdrawHandler).Methods("POST")
	bankingRouter.HandleFunc("/account/transfer", r.transferHandler).Methods("POST")
}

// --------- API struct TYPES -----------

type createAccountRequest struct {
	Currency string `json:"currency"`
}

type depositWithdrawRequest struct {
	Amount float64 `json:"amount"`
}

type transferRequest struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

// ----------- HANDLERS ------------

func (r *Router) createAccountHandler(w http.ResponseWriter, req *http.Request) {
	userID, err := middleware.ValidateUser(req)
	if err != nil {
		r.logger.WithError(err).Error("failed to authenticate user")
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	var reqBody createAccountRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	account, err := r.bankingService.CreateAccount(req.Context(), userID, reqBody.Currency)
	if err != nil {
		r.logger.WithError(err).Error("account creation fail")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (r *Router) depositHandler(w http.ResponseWriter, req *http.Request) {
	userID, err := middleware.ValidateUser(req)
	if err != nil {
		r.logger.WithError(err).Error("failed to authenticate user")
		http.Error(w, "Invalid request", http.StatusUnauthorized)
		return
	}
	accountID, err := middleware.ValidateAccount(req)
	if err != nil {
		http.Error(w, "Invalid account", http.StatusBadRequest)
		return
	}
	var reqBody depositWithdrawRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil || reqBody.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	account, err := r.bankingService.GetAccountByID(context.Background(), accountID)
	if err != nil || account.UserID != userID {
		r.logger.WithError(err).Error("failed to get account")
		http.Error(w, "could not get account", http.StatusInternalServerError)
		return
	}
	if err := r.bankingService.Deposit(req.Context(), account.ID, reqBody.Amount); err != nil {
		r.logger.WithError(err).Error("failed to deposit")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (r *Router) withdrawHandler(w http.ResponseWriter, req *http.Request) {
	userID, err := middleware.ValidateUser(req)
	if err != nil {
		r.logger.WithError(err).Error("failed to authenticate user")
		http.Error(w, "Invalid request", http.StatusUnauthorized)
		return
	}
	accountID, err := middleware.ValidateAccount(req)
	if err != nil {
		http.Error(w, "Invalid account", http.StatusBadRequest)
		return
	}
	var reqBody depositWithdrawRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil || reqBody.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	account, err := r.bankingService.GetAccountByID(context.Background(), accountID)
	if err != nil || account.UserID != userID {
		r.logger.WithError(err).Error("failed to get account")
		http.Error(w, "could not get account", http.StatusInternalServerError)
		return
	}
	if err := r.bankingService.Withdraw(req.Context(), accountID, reqBody.Amount); err != nil {
		r.logger.WithError(err).Error("failed to withdraw")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (r *Router) transferHandler(w http.ResponseWriter, req *http.Request) {
	userID, err := middleware.ValidateUser(req)
	if err != nil {
		r.logger.WithError(err).Error("failed to authenticate user")
		http.Error(w, "Invalid request", http.StatusUnauthorized)
		return
	}
	var reqBody transferRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if reqBody.FromAccountID == reqBody.ToAccountID || reqBody.Amount <= 0 {
		http.Error(w, "Invalid transfer parameters", http.StatusBadRequest)
		return
	}
	account, err := r.bankingService.GetAccountByID(context.Background(), reqBody.FromAccountID)
	if err != nil || account.UserID != userID {
		r.logger.WithError(err).Error("failed to get account")
		http.Error(w, "Could not get account", http.StatusBadRequest)
		return
	}
	if err := r.bankingService.Transfer(req.Context(), reqBody.FromAccountID, reqBody.ToAccountID, reqBody.Amount); err != nil {
		r.logger.WithError(err).Error("Transfer failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// func (r *Router)  getAccountByIDHandler(w http.ResponseWriter, req *http.Request) {
// 	accountID, err := parseIDFromVars(r, "id")
// 	if err != nil {
// 		http.Error(w, "Invalid account id", http.StatusBadRequest)
// 		return
// 	}
// 	account, err := r.bankingService.GetAccountByID(r.Context(), accountID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(account)
// }

// func (r *Router)  getAccountsByUserHandler(w http.ResponseWriter, req *http.Request) {
// 	var (
// 		userID string
// 		ok     bool
// 	)
// 	if userID, ok = r.Context().Value(middleware.UserIDKey).(string); !ok {
// 		http.Error(w, "Invalid user", http.StatusUnauthorized)
// 		return
// 	}
// 	accounts, err := r.bankingService.GetAccountsByUser(r.Context(), userID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		r.logger.WithError(err).Error("getAccountsByUser failed")
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(accounts)
// }

// ----------- HELPERS --------------

func parseIDFromVars(req *http.Request, varName string) (int64, error) {
	vars := mux.Vars(req)
	raw, ok := vars[varName]
	if !ok {
		return 0, errors.New("missing id")
	}
	return strconv.ParseInt(raw, 10, 64)
}
