package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mshagirov/chirpy/internal/auth"
	"github.com/mshagirov/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) serveCreateUser(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		errorResponse(fmt.Sprintf("Error decoding email %s and password", err),
			http.StatusBadRequest, w, r)
		return
	}
	if !validEmail(params.Email) {
		errorResponse(fmt.Sprintf("Not a valid email (%s)", params.Email),
			http.StatusBadRequest, w, r)
		return
	}
	if len(params.Password) < 1 {
		errorResponse("Empty password field", http.StatusBadRequest, w, r)
	}
	passHash, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResponse("Error saving password", http.StatusServiceUnavailable, w, r)
	}
	user, err := cfg.db.CreateUser(r.Context(),
		database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: passHash,
		})
	if err != nil {
		errorResponse(fmt.Sprintf("Database error: %s", err),
			http.StatusServiceUnavailable, w, r)
		return
	}
	jsonResponse(User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, http.StatusCreated, w, r)
}

func validEmail(email string) bool {
	if !strings.ContainsRune(email, '@') {
		return false
	}
	if strings.ContainsRune(email, ' ') {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) < 1 {
		return false
	}
	domain := strings.Split(parts[1], ".")
	if len(domain) != 2 {
		return false
	}
	return true
}
