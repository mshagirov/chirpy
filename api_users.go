package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) serveCreateUser(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Email string `json:"email"`
	}{}
	var err_str string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		err_str = fmt.Sprintf("Error decoding email: %s", err)
		log.Println(err_str)
		errorResponse(err_str, http.StatusBadRequest, w, r)
		return
	}
	if !validEmail(params.Email) {
		err_str = fmt.Sprintf("Not a valid email: %s", params.Email)
		log.Println(err_str)
		errorResponse(err_str, http.StatusBadRequest, w, r)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		err_str = fmt.Sprintf("Database error: %s", err)
		log.Println(err_str)
		errorResponse(err_str, http.StatusServiceUnavailable, w, r)
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
