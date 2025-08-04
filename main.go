package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/mshagirov/chirpy/internal/database"
)

func main() {
	const rootPath = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	// CLI command for rand num: openssl rand -base64 64
	secret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{}
	apiCfg.db = dbQueries
	apiCfg.platform = platform
	apiCfg.secret = secret
	apiCfg.fileserverHits.Store(0)

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.serveMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.serveReset)
	mux.HandleFunc("GET /api/healthz", ServerReady)
	mux.HandleFunc("POST /api/users", apiCfg.serveCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.serveLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.serveRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.serveRevoke)
	mux.HandleFunc("POST /api/chirps", apiCfg.serveCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.serveGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.serveGetChirpWithID)
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s port: %s\n",
		rootPath,
		port)
	log.Fatal(srv.ListenAndServe())
}
