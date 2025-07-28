package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) serveReset(w http.ResponseWriter, r *http.Request) {
	var err_str string
	if cfg.platform != "dev" {
		err_str = fmt.Sprintf("Forbidden: platform=%s", cfg.platform)
		log.Println(err_str)
		errorResponse(err_str, http.StatusForbidden, w, r)
		return
	}
	// delete ALL users
	err := cfg.db.Reset(r.Context())
	if err != nil {
		err_str = fmt.Sprintf("Error reseting database: %s", err)
		log.Println(err_str)
		errorResponse(err_str, http.StatusServiceUnavailable, w, r)
	}
	log.Println("Reset database; deleted ALL users")
	// reset visit counter
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Counter was reset to 0")); err != nil {
		log.Println("apiConfig.serveReset error: ", err)
	}
	log.Println("Reset fileserverHits to ", cfg.fileserverHits.Load())
}
