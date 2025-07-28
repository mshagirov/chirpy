package main

import (
	"log"
	"net/http"
)

func ServerReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Status : 200 " + http.StatusText(http.StatusOK)))
	if err != nil {
		log.Printf("ServerReady handler got error: %v", err)
	}
}
