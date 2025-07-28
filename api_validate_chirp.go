package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func isBadWord(s string) bool {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	return slices.Contains(badWords, strings.ToLower(s))
}

func serveValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpParameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := chirpParameters{}
	if err := decoder.Decode(&params); err != nil {
		errorResponse(
			fmt.Sprintf("Error decoding parameters: %s", err),
			http.StatusInternalServerError, w, r)
		return
	}
	if len([]rune(params.Body)) > 140 {
		errorResponse("Chirp is too long",
			http.StatusBadRequest, w, r)
		return
	}
	words := strings.Split(params.Body, " ")
	for idx, val := range words {
		if isBadWord(val) {
			words[idx] = "****"
		}
	}
	params.Body = strings.Join(words, " ")
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	jsonResponse(returnVals{CleanedBody: params.Body}, http.StatusOK, w, r)
}
