package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
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

func censorBadWords(s string) string {
	isBadWord := func(w string) bool {
		badWords := []string{"kerfuffle", "sharbert", "fornax"}
		return slices.Contains(badWords, strings.ToLower(w))
	}
	words := strings.Split(s, " ")
	for idx, val := range words {
		if isBadWord(val) {
			words[idx] = "****"
		}
	}
	return strings.Join(words, " ")
}
