package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mshagirov/chirpy/internal/auth"
	"github.com/mshagirov/chirpy/internal/database"
)

func (cfg *apiConfig) serveLogin(w http.ResponseWriter, r *http.Request) {
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
	user, err := cfg.db.GetUserWithEmail(r.Context(), params.Email)
	if err != nil || auth.CheckPasswordHash(params.Password, user.HashedPassword) != nil {
		log.Println("api/login error:", err)
		errorResponse("Incorrect email or password", http.StatusUnauthorized, w, r)
		return
	}
	randomString, err := auth.MakeRefreshToken()
	if err != nil {
		errorResponse(fmt.Sprintf("Error creating refresh token: %v", err),
			http.StatusBadRequest, w, r)
		return
	}
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     randomString,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 1440),
	})
	if err != nil {
		errorResponse(fmt.Sprintf("Database error creating refresh token: %v", err),
			http.StatusBadRequest, w, r)
		return
	}
	tokenString, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		errorResponse("Error creating token", http.StatusUnauthorized, w, r)
	}
	jsonResponse(User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: refreshToken.Token,
	}, http.StatusOK, w, r)
}
