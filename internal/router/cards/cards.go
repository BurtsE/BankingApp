package cards

import (
	"BankingApp/internal/service"
	"net/http"

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
	secured := c.muxRouter.NewRoute().Subrouter()
	// secured.Use(router.jwtMiddleware) // JWT middleware

	secured.HandleFunc("/issue", c.issueCardHandler).Methods("POST")
	secured.HandleFunc("/show", c.showCardHandler).Methods("GET")
	return c
}

func (c *CardSubRouter) issueCardHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Карта выпущена"}`))
}
func (c *CardSubRouter) showCardHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Описание карты"}`))
}
