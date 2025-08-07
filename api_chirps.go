package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mshagirov/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) serveCreateChirp(w http.ResponseWriter, r *http.Request) {
	user_id, err := jwtGetUserID(cfg.secret, w, r)
	if err != nil {
		return
	}
	params := struct {
		Body string `json:"body"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		errorResponse(fmt.Sprintf("Error decoding parameters: %s", err), http.StatusInternalServerError, w, r)
		return
	}
	if len([]rune(params.Body)) > 140 {
		errorResponse("Chirp is too long", http.StatusBadRequest, w, r)
		return
	}
	params.Body = censorBadWords(params.Body)
	type CreateChirpParams struct {
		Body   string
		UserID uuid.UUID
	}
	chirpParams, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: user_id,
	})
	if err != nil {
		errorResponse(fmt.Sprintf("Database error: %s", err), http.StatusServiceUnavailable, w, r)
		return
	}
	jsonResponse(Chirp{
		ID:        chirpParams.ID,
		CreatedAt: chirpParams.CreatedAt,
		UpdatedAt: chirpParams.UpdatedAt,
		Body:      chirpParams.Body,
		UserID:    chirpParams.UserID,
	}, http.StatusCreated, w, r)
}

func (cfg *apiConfig) serveGetChirps(w http.ResponseWriter, r *http.Request) {
	var user_id uuid.UUID
	var chirps []database.Chirp
	var err error
	isForAuthor := false
	// Query: api/chirps?author_id=AUTHORS
	if authorIDStr := r.URL.Query().Get("author_id"); len(authorIDStr) > 0 {
		user_id, err = uuid.Parse(authorIDStr)
		if err == nil {
			isForAuthor = true
		}
	}
	if isForAuthor {
		chirps, err = cfg.db.GetChirpsWithUserID(r.Context(), user_id)
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
	}
	if err != nil {
		errorResponse(fmt.Sprintf("Database error: %s", err), http.StatusServiceUnavailable, w, r)
		return
	}
	var returnVals []Chirp
	for _, c := range chirps {
		returnVals = append(returnVals,
			Chirp{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserID:    c.UserID,
			})
	}
	jsonResponse(returnVals, http.StatusOK, w, r)
}

func (cfg *apiConfig) serveGetChirpWithID(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(strID)
	if err != nil {
		errorResponse(fmt.Sprintf("Error parsing chirpID: %s", err), http.StatusNotFound, w, r)
		return
	}
	c, err := cfg.db.GetChirpWithID(r.Context(), chirpID)
	if err != nil {
		errorResponse("Chirp not found", http.StatusNotFound, w, r)
		return
	}
	jsonResponse(Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}, http.StatusOK, w, r)
}

func (cfg *apiConfig) serveDeleteChirpWithID(w http.ResponseWriter, r *http.Request) {
	user_id, err := jwtGetUserID(cfg.secret, w, r)
	if err != nil {
		return
	}
	strID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(strID)
	if err != nil {
		errorResponse(fmt.Sprintf("Error parsing chirpID: %s", err), http.StatusNotFound, w, r)
		return
	}
	c, err := cfg.db.GetChirpWithID(r.Context(), chirpID)
	if err != nil {
		errorResponse("Chirp not found", http.StatusNotFound, w, r)
		return
	}
	if c.UserID != user_id {
		errorResponse("Forbidden", http.StatusForbidden, w, r)
		return
	}
	err = cfg.db.DeleteChirpWithID(r.Context(), c.ID)
	if err != nil {
		errorResponse(fmt.Sprintf("Database error: %v", err), http.StatusBadRequest, w, r)
	}
	w.WriteHeader(http.StatusNoContent)
}
