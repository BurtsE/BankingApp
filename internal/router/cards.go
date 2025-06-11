package router

import (
	"BankingApp/internal/config"
	"BankingApp/pkg/middleware"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// --- PROTECTED ROUTES (JWT Auth Required) ---

func (r *Router) InitCardRoutes() {

	authMiddleware := middleware.NewAuthMiddleware(config.GetJWTSecretKey())
	cardRouter := r.muxRouter.PathPrefix("/card").Subrouter()
	cardRouter.Use(authMiddleware)
	cardRouter.HandleFunc("/issue", r.issueCardHandler).Methods("POST")
	cardRouter.HandleFunc("/show", r.showCardHandler).Methods("GET")
}

// --------- API struct TYPES -----------

type issueCardRequest struct {
	AccountId int64 `json:"account_id"`
}

type showCardsRequest struct {
	AccountId int64 `json:"account_id"`
}

// ----------- HANDLERS ------------

func (r *Router) issueCardHandler(w http.ResponseWriter, req *http.Request) {
	UUID, err := middleware.ValidateUser(req)
	if err != nil {
		r.logger.WithError(err).Error("failed to authenticate user")
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	var reqBody issueCardRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		r.logger.Println(err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	user, err := r.userService.GetByID(ctx, UUID)
	if err != nil {
		r.logger.WithError(err).Error("failed to get user")
		http.Error(w, "could not get user", http.StatusInternalServerError)
		return
	}
	account, err := r.bankingService.GetAccountByID(ctx, reqBody.AccountId)
	if err != nil || account.UserID != UUID {
		r.logger.WithError(err).Error("failed to get account")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	card, err := r.cardService.GenerateVirtualCard(ctx, reqBody.AccountId, user.FullName)
	if err != nil {
		r.logger.WithError(err).Error("failed to generate card")
		http.Error(w, "card generation error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}
func (r *Router) showCardHandler(w http.ResponseWriter, req *http.Request) {
	UUID, err := middleware.ValidateUser(req)
	if err != nil {
		r.logger.WithError(err).Error("failed to authenticate user")
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	var reqBody showCardsRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	account, err := r.bankingService.GetAccountByID(ctx, reqBody.AccountId)
	if err != nil || account.UserID != UUID {
		r.logger.WithError(err).Error("failed to get account")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	cards, err := r.cardService.GetCardsByAccount(req.Context(), reqBody.AccountId)
	if err != nil || account.UserID != UUID {
		r.logger.WithError(err).Errorf("failed to get cards for account: %w", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)

}
