package main

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/mshagirov/chirpy/internal/auth"
)

func jwtGetUserID(secret string, w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	var user_id uuid.UUID
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse("Unauthorized", http.StatusUnauthorized, w, r)
		return user_id, err
	}
	user_id, err = auth.ValidateJWT(bearerToken, secret)
	if err != nil {
		errorResponse("Unauthorized", http.StatusUnauthorized, w, r)
		return user_id, err
	}
	return user_id, nil
}
