package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/mshagirov/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) serveMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	htmlBody := fmt.Sprintf(`<html>
 <body>
	<h1>Welcome, Chirpy Admin</h1>
	 <p>Chirpy has been visited %d times!</p>
 </body>
</html>`, cfg.fileserverHits.Load())
	_, err := w.Write([]byte(htmlBody))
	if err != nil {
		log.Printf("apiConfig.serveMetrics got error: %v", err)
	}
}

func (cfg *apiConfig) serveReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Counter was reset to 0")); err != nil {
		log.Println("apiConfig.serveReset error: ", err)
	}
	log.Println("Reset fileserverHits to ", cfg.fileserverHits.Load())
}
