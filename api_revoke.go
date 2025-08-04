package main

import (
	"net/http"

	"github.com/mshagirov/chirpy/internal/auth"
)

func (cfg *apiConfig) serveRevoke(w http.ResponseWriter, r *http.Request) {
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

	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		errorResponse("Unauthorized", http.StatusUnauthorized, w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
