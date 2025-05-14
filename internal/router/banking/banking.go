package banking

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"BankingApp/internal/service"

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
		muxRouter:      muxRouter,
		logger:         logger,
		bankingService: bankingService,
	}
	br.routes()
	return br
}

func (br *BankingSubRouter) routes() {
	br.muxRouter.HandleFunc("/account", br.createAccountHandler).Methods("POST")
	br.muxRouter.HandleFunc("/account/{id:[0-9]+}/deposit", br.depositHandler).Methods("POST")
	br.muxRouter.HandleFunc("/account/{id:[0-9]+}/withdraw", br.withdrawHandler).Methods("POST")
	br.muxRouter.HandleFunc("/account/transfer", br.transferHandler).Methods("POST")
	br.muxRouter.HandleFunc("/account/{id:[0-9]+}", br.getAccountByIDHandler).Methods("GET")
	br.muxRouter.HandleFunc("/user/{id:[0-9]+}/accounts", br.getAccountsByUserHandler).Methods("GET")
}

// --------- API struct TYPES -----------

type createAccountRequest struct {
	UserID   int64  `json:"user_id"`
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
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	account, err := br.bankingService.CreateAccount(r.Context(), req.UserID, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		br.logger.WithError(err).Error("CreateAccount failed")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (br *BankingSubRouter) depositHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := parseIDFromVars(r, "id")
	if err != nil {
		http.Error(w, "Invalid account id", http.StatusBadRequest)
		return
	}
	var req depositWithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	if err := br.bankingService.Deposit(r.Context(), accountID, req.Amount); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		br.logger.WithError(err).Error("Deposit failed")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (br *BankingSubRouter) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := parseIDFromVars(r, "id")
	if err != nil {
		http.Error(w, "Invalid account id", http.StatusBadRequest)
		return
	}
	var req depositWithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
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
	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.FromAccountID == req.ToAccountID || req.Amount <= 0 {
		http.Error(w, "Invalid transfer parameters", http.StatusBadRequest)
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

func (br *BankingSubRouter) getAccountByIDHandler(w http.ResponseWriter, r *http.Request) {
	accountID, err := parseIDFromVars(r, "id")
	if err != nil {
		http.Error(w, "Invalid account id", http.StatusBadRequest)
		return
	}
	account, err := br.bankingService.GetAccountByID(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (br *BankingSubRouter) getAccountsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := parseIDFromVars(r, "id")
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	accounts, err := br.bankingService.GetAccountsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		br.logger.WithError(err).Error("getAccountsByUser failed")
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}

// ----------- HELPERS --------------

func parseIDFromVars(r *http.Request, varName string) (int64, error) {
	vars := mux.Vars(r)
	raw, ok := vars[varName]
	if !ok {
		return 0, errors.New("missing id")
	}
	return strconv.ParseInt(raw, 10, 64)
}
