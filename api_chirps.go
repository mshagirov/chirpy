package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	var err_str string
	params := struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		err_str = fmt.Sprintf("Error decoding parameters: %s", err)
		log.Println(err_str)
		errorResponse(err_str, http.StatusInternalServerError, w, r)
		return
	}
	if len([]rune(params.Body)) > 140 {
		errorResponse("Chirp is too long", http.StatusBadRequest, w, r)
		return
	}
	user_id, err := uuid.Parse(params.UserID)
	if err != nil {
		err_str = fmt.Sprintf("Error parsing user_id: %s", err)
		log.Println(err_str)
		errorResponse(err_str, http.StatusBadRequest, w, r)
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
		err_str = fmt.Sprintf("Database error: %s", err)
		log.Println(err_str)
		errorResponse(err_str, http.StatusServiceUnavailable, w, r)
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
	var err_str string
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		err_str = fmt.Sprintf("Database error: %s", err)
		log.Println(err_str)
		errorResponse(err_str, http.StatusServiceUnavailable, w, r)
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
	var err_str string
	strID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(strID)
	if err != nil {
		err_str = fmt.Sprintf("Error parsing chirpID: %s", err)
		log.Println(err_str)
		errorResponse(err_str, http.StatusNotFound, w, r)
		return
	}
	c, err := cfg.db.GetChirpWithID(r.Context(), chirpID)
	if err != nil {
		err_str = "Chirp not found!"
		log.Println(err_str)
		errorResponse(err_str, http.StatusNotFound, w, r)
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
