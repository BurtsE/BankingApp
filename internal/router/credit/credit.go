package credit

import (
	"BankingApp/internal/service"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type CreditSubRouter struct {
	muxRouter     *mux.Router
	logger        *logrus.Logger
	creditService service.CreditService
}

func InitCardRouter(creditService service.CreditService, logger *logrus.Logger, muxRouter *mux.Router) *CreditSubRouter {
	c := &CreditSubRouter{
		muxRouter:   muxRouter,
		logger:      logger,
		creditService: creditService,
	}
	secured := c.muxRouter.NewRoute().Subrouter()
	// secured.Use(router.jwtMiddleware) // JWT middleware

	secured.HandleFunc("/issue", c.issueCreditHandler).Methods("POST")
	secured.HandleFunc("/schedule", c.showCardHandler).Methods("POST")
	secured.HandleFunc("/payment-graph", c.showPaymentsHandler).Methods("POST")
	return c
}

func (c *CreditSubRouter) issueCreditHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Оформление кредита с расчетом аннуитетных платежей"}`))
}
func (c *CreditSubRouter) showCardHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Автоматическое списание платежей (шедулер каждые N часов)"}`))
}

func (c *CreditSubRouter) showPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Генерация графика платежей"}`))
}
