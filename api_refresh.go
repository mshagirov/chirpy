package main

import (
	"net/http"
	"time"

	"github.com/mshagirov/chirpy/internal/auth"
)

func (cfg *apiConfig) serveRefresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse("Unauthorized", http.StatusUnauthorized, w, r)
		return
	}
	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		errorResponse("Unauthorized", http.StatusUnauthorized, w, r)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		errorResponse("Unauthorized", http.StatusUnauthorized, w, r)
		return
	}
	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		errorResponse("Error creating token", http.StatusUnauthorized, w, r)
		return
	}

	jsonResponse(struct {
		Token string `json:"token"`
	}{Token: newToken}, http.StatusOK, w, r)
}
