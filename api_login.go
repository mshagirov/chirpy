package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mshagirov/chirpy/internal/auth"
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
		if err != nil {
			log.Println("api/login error:", err)
		}
		errorResponse("Incorrect email or password", http.StatusUnauthorized, w, r)
		return
	}
	jsonResponse(User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}, http.StatusOK, w, r)
}
