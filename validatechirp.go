package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func errorResponse(msg string, code int, w http.ResponseWriter, r *http.Request) {
	type errorVal struct {
		Error string `json:"error"`
	}
	jsonResponse(errorVal{Error: msg}, code, w, r)
}

func jsonResponse(payload any, code int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("%d: Error marshalling JSON: %s", http.StatusInternalServerError, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("%d: %v %v %v %v", code, r.Method, r.Proto, r.URL.String(), r.Header.Get("User-Agent"))
	w.WriteHeader(code)
	w.Write(data)
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
	type returnVals struct {
		Valid bool `json:"valid"`
	}
	jsonResponse(returnVals{Valid: true}, http.StatusOK, w, r)
}
