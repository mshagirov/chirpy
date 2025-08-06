package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"github.com/mshagirov/chirpy/internal/auth"
	"github.com/mshagirov/chirpy/internal/database"
)

func (cfg *apiConfig) servePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	params := struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		errorResponse(fmt.Sprintf("Error decoding email %s and password", err),
			http.StatusBadRequest, w, r)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	user_id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		errorResponse(fmt.Sprintf("Error parsing user_id: %s", err), http.StatusNotFound, w, r)
		return
	}
	chirpyRedParams := database.SetChirpyRedWithIDParams{
		ID:          user_id,
		IsChirpyRed: true,
	}
	user, err := cfg.db.SetChirpyRedWithID(r.Context(), chirpyRedParams)
	if err != nil || !user.IsChirpyRed {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
