package cards

import (
	"BankingApp/internal/config"
	"BankingApp/internal/service"
	"BankingApp/pkg/middleware"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// --- PROTECTED ROUTES (JWT Auth Required) ---
type CardSubRouter struct {
	muxRouter   *mux.Router
	logger      *logrus.Logger
	cardService service.CardService
}

func InitCardRouter(cardService service.CardService, logger *logrus.Logger, muxRouter *mux.Router) *CardSubRouter {
	c := &CardSubRouter{
		muxRouter:   muxRouter.PathPrefix("/cards").Subrouter(),
		logger:      logger,
		cardService: cardService,
	}
	authMiddleware := middleware.NewAuthMiddleware(config.GetJWTSecretKey())
	c.muxRouter.Use(authMiddleware)

	c.routes()
	return c
}

func (c *CardSubRouter) routes() {
	c.muxRouter.HandleFunc("/issue", c.issueCardHandler).Methods("POST")
	c.muxRouter.HandleFunc("/show", c.showCardHandler).Methods("GET")
}

// --------- API struct TYPES -----------

type issueCardRequest struct {
	AccountId int64  `json:"account_id"`
	UserName  string `json:"user_name`
}

type showCardsRequest struct {
	accountId int64 `json:"account_id"`
}

// ----------- HANDLERS ------------

func (c *CardSubRouter) issueCardHandler(w http.ResponseWriter, r *http.Request) {
	_, err := middleware.ValidateUser(r)
	if err != nil {
		http.Error(w, "Invalid user", http.StatusUnauthorized)
		return
	}
	var req issueCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	card, err := c.cardService.GenerateVirtualCard(ctx, req.AccountId, req.UserName)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}
func (c *CardSubRouter) showCardHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Описание карты"}`))
}
