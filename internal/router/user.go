package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (r *Router) InitUserRoutes() {
	userRouter := r.muxRouter.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", r.registerUserHandler).Methods("POST")
	userRouter.HandleFunc("/login", r.loginHandler).Methods("POST")
	userRouter.HandleFunc("/{id:[0-9]+}", r.getUserByIDHandler).Methods("GET")
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

func (r *Router) registerUserHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody registerRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	user, err := r.userService.Register(ctx, reqBody.Email, reqBody.Username, reqBody.Password, reqBody.FullName)
	if err != nil {
		r.logger.WithError(err).Error("failed to register user")
		http.Error(w, "Registration failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (r *Router) loginHandler(w http.ResponseWriter, req *http.Request) {
	var reqBody loginRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	jwtToken, expiresAt, err := r.userService.Authenticate(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		r.logger.WithError(err).Warn("failed to authenticate user")
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":      jwtToken,
		"expires_at": expiresAt,
	})
}

func (r *Router) getUserByIDHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID, ok := vars["id"]
	if !ok {
		http.Error(w, "User ID not specified", http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	user, err := r.userService.GetByID(ctx, userID)
	if err != nil {
		r.logger.WithError(err).Warn("user not found")
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
