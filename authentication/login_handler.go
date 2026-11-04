package authentication

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type handlerResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type LoginHandler struct {
	service Authentication
}

func NewLoginHandler(s Authentication) *LoginHandler {
	return &LoginHandler{service: s}
}

func (l *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req UserReq
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	if err := l.service.Login(req.Username, req.Password); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	writeJSON(w, http.StatusOK, handlerResponse{Message: "logged in"})
}

func (l *LoginHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserReq
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	user, err := l.service.CreateUser(&req)
	switch {
	case errors.Is(err, ErrInvalidPassword):
		writeError(w, http.StatusBadRequest, "password is too short")
		return
	case err != nil:
		log.Printf("create user: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, handlerResponse{Message: "user created", Data: user.ID})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, handlerResponse{Message: msg})
}
