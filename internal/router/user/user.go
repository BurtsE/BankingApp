package user

import (
	"BankingApp/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UserSubRouter struct {
	muxRouter   *mux.Router
	logger      *logrus.Logger
	userService service.UserService
}

func InitUserRouter(userService service.UserService, logger *logrus.Logger, muxRouter *mux.Router) *UserSubRouter {
	u := &UserSubRouter{
		muxRouter:   muxRouter.PathPrefix("/user").Subrouter(),
		logger:      logger,
		userService: userService,
	}
	u.muxRouter.HandleFunc("/register", u.registerHandler).Methods("POST")
	u.muxRouter.HandleFunc("/login", u.loginHandler).Methods("POST")
	u.muxRouter.HandleFunc("/{id:[0-9]+}", u.getUserByIDHandler).Methods("GET")
	return u
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserSubRouter) registerHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user, err := u.userService.Register(ctx, req.Email, req.Username, req.Password, req.FullName)
	if err != nil {
		u.logger.WithError(err).Error("failed to register user")
		http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (u *UserSubRouter) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	jwtToken, expiresAt, err := u.userService.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		u.logger.WithError(err).Warn("failed to authenticate user")
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      jwtToken,
		"expires_at": expiresAt,
	})
}

func (u *UserSubRouter) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID not specified", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	user, err := u.userService.GetByID(ctx, userID)
	if err != nil {
		u.logger.WithError(err).Warn("user not found")
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
