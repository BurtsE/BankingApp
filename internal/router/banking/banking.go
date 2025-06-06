package banking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"BankingApp/internal/config"
	"BankingApp/internal/service"
	"BankingApp/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// --- PROTECTED ROUTES (JWT Auth Required) ---

type BankingSubRouter struct {
	muxRouter      *mux.Router
	logger         *logrus.Logger
	bankingService service.BankingService
}

func InitBankingRouter(bankingService service.BankingService, logger *logrus.Logger, muxRouter *mux.Router) *BankingSubRouter {
	br := &BankingSubRouter{
		muxRouter:      muxRouter.PathPrefix("/banking").Subrouter(),
		logger:         logger,
		bankingService: bankingService,
	}
	authMiddleware := middleware.NewAuthMiddleware(config.GetJWTSecretKey())
	br.muxRouter.Use(authMiddleware)
	br.routes()
	return br
}

func (br *BankingSubRouter) routes() {
	br.muxRouter.HandleFunc("/account", br.createAccountHandler).Methods("POST")
	br.muxRouter.HandleFunc("/account/{id:[0-9]+}/deposit", br.depositHandler).Methods("POST")
	br.muxRouter.HandleFunc("/account/{id:[0-9]+}withdraw", br.withdrawHandler).Methods("POST")
	br.muxRouter.HandleFunc("/account/transfer", br.transferHandler).Methods("POST")
	// br.muxRouter.HandleFunc("/account/{id:[0-9]+}", br.getAccountByIDHandler).Methods("GET")
	// br.muxRouter.HandleFunc("/account/user/{id:[0-9]+}/accounts", br.getAccountsByUserHandler).Methods("GET")
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

func (br *BankingSubRouter) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := validateUser(r)
	if err != nil {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	account, err := br.bankingService.CreateAccount(r.Context(), userID, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		br.logger.WithError(err).Error("CreateAccount failed")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (br *BankingSubRouter) depositHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := validateUser(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusUnauthorized)
		return
	}
	accountID, err := validateAccount(r)
	if err != nil {
		http.Error(w, "Invalid account", http.StatusBadRequest)
		return
	}
	var req depositWithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	account, err := br.bankingService.GetAccountByID(context.Background(), accountID)
	if err != nil || account.UserID != userID {
		http.Error(w, "could not get account", http.StatusInternalServerError)
		return
	}
	if err := br.bankingService.Deposit(r.Context(), account.ID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		br.logger.WithError(err).Error("Deposit failed")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (br *BankingSubRouter) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := validateUser(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusUnauthorized)
		return
	}
	accountID, err := validateAccount(r)
	if err != nil {
		http.Error(w, "Invalid account", http.StatusBadRequest)
		return
	}
	var req depositWithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	account, err := br.bankingService.GetAccountByID(context.Background(), accountID)
	if err != nil || account.UserID != userID {
		http.Error(w, "could not get account", http.StatusInternalServerError)
		return
	}
	if err := br.bankingService.Withdraw(r.Context(), accountID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		br.logger.WithError(err).Error("Withdraw failed")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (br *BankingSubRouter) transferHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := validateUser(r)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusUnauthorized)
		return
	}
	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.FromAccountID == req.ToAccountID || req.Amount <= 0 {
		http.Error(w, "Invalid transfer parameters", http.StatusBadRequest)
		return
	}
	account, err := br.bankingService.GetAccountByID(context.Background(), req.FromAccountID)
	if err != nil || account.UserID != userID {
		http.Error(w, "Could not get account", http.StatusBadRequest)
		return
	}
	if err := br.bankingService.Transfer(r.Context(), req.FromAccountID, req.ToAccountID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		br.logger.WithError(err).Error("Transfer failed")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// func (br *BankingSubRouter) getAccountByIDHandler(w http.ResponseWriter, r *http.Request) {
// 	accountID, err := parseIDFromVars(r, "id")
// 	if err != nil {
// 		http.Error(w, "Invalid account id", http.StatusBadRequest)
// 		return
// 	}
// 	account, err := br.bankingService.GetAccountByID(r.Context(), accountID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(account)
// }

// func (br *BankingSubRouter) getAccountsByUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		userID string
// 		ok     bool
// 	)
// 	if userID, ok = r.Context().Value(middleware.UserIDKey).(string); !ok {
// 		http.Error(w, "Invalid user", http.StatusUnauthorized)
// 		return
// 	}
// 	accounts, err := br.bankingService.GetAccountsByUser(r.Context(), userID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		br.logger.WithError(err).Error("getAccountsByUser failed")
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(accounts)
// }

// ----------- HELPERS --------------

func validateUser(r *http.Request) (string, error) {
	var (
		userID string
		ok     bool
	)
	if userID, ok = r.Context().Value(middleware.UserIDKey).(string); !ok {
		// http.Error(w, "Invalid user", http.StatusUnauthorized)
		return "", fmt.Errorf("not authenticated user")
	}

	return userID, nil
}

func validateAccount(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	accountIDStr := vars["id"]
	if accountIDStr == "" {
		// http.Error(w, "could not get account", http.StatusInternalServerError)
		return 0, fmt.Errorf(" missing account id")
	}
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		// http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return 0, fmt.Errorf(" missing account id")
	}
	return accountID, nil
}

func parseIDFromVars(r *http.Request, varName string) (int64, error) {
	vars := mux.Vars(r)
	raw, ok := vars[varName]
	if !ok {
		return 0, errors.New("missing id")
	}
	return strconv.ParseInt(raw, 10, 64)
}
