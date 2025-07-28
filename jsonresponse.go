package main

import (
	"encoding/json"
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
