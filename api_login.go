package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mshagirov/chirpy/internal/auth"
)

func (cfg *apiConfig) serveLogin(w http.ResponseWriter, r *http.Request) {
	var duration time.Duration
	params := struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		errorResponse(fmt.Sprintf("Error decoding email %s and password", err),
			http.StatusBadRequest, w, r)
		return
	}
	if params.ExpiresInSeconds < 1 || params.ExpiresInSeconds > 3600 {
		duration = 3600 * time.Second
	} else {
		duration = time.Duration(params.ExpiresInSeconds) * time.Second
	}
	user, err := cfg.db.GetUserWithEmail(r.Context(), params.Email)
	if err != nil || auth.CheckPasswordHash(params.Password, user.HashedPassword) != nil {
		log.Println("api/login error:", err)
		errorResponse("Incorrect email or password", http.StatusUnauthorized, w, r)
		return
	}
	tokenString, err := auth.MakeJWT(user.ID, cfg.secret, duration)
	if err != nil {
		errorResponse("Error creating token", http.StatusUnauthorized, w, r)
	}
	jsonResponse(User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     tokenString,
	}, http.StatusOK, w, r)
}
