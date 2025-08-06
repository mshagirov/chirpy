package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mshagirov/chirpy/internal/auth"
	"github.com/mshagirov/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type userRequestParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func decodeUserRequest(w http.ResponseWriter, r *http.Request) (userRequestParams, error) {
	params := userRequestParams{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		errorResponse(fmt.Sprintf("Error decoding email and password: %v", err),
			http.StatusBadRequest, w, r)
		return params, err
	}
	if !validEmail(params.Email) {
		errorResponse(fmt.Sprintf("Invalid email (%s)", params.Email),
			http.StatusBadRequest, w, r)
		return params, fmt.Errorf("Invalid email")
	}
	if len(params.Password) < 1 {
		errorResponse("Empty password field", http.StatusBadRequest, w, r)
		return params, fmt.Errorf("Empty password")
	}
	return params, nil
}

func (cfg *apiConfig) serveCreateUser(w http.ResponseWriter, r *http.Request) {
	params, err := decodeUserRequest(w, r)
	if err != nil {
		return
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
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, http.StatusCreated, w, r)
}

func (cfg *apiConfig) serveUpdateUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := jwtGetUserID(cfg.secret, w, r)
	if err != nil {
		return
	}
	params, err := decodeUserRequest(w, r)
	if err != nil {
		return
	}
	passHash, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResponse("Error saving password", http.StatusServiceUnavailable, w, r)
	}
	updateParams := database.UpdateUserWithIDParams{
		ID:             user_id,
		Email:          params.Email,
		HashedPassword: passHash,
	}
	user, err := cfg.db.UpdateUserWithID(r.Context(), updateParams)
	if err != nil {
		errorResponse(fmt.Sprintf("Database error: %s", err),
			http.StatusServiceUnavailable, w, r)
		return
	}
	jsonResponse(User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, http.StatusOK, w, r)
}
