package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
