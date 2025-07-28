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
	apiCfg := apiConfig{}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	apiCfg.db = dbQueries
	const rootPath = "."
	const port = "8080"
	apiCfg.fileserverHits.Store(0)
	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", ServerReady)
	mux.HandleFunc("GET /admin/metrics", apiCfg.serveMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.serveReset)
	mux.HandleFunc("POST /api/validate_chirp", serveValidateChirp)
	mux.HandleFunc("POST /api/users", apiCfg.serveCreateUser)
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s port: %s\n",
		rootPath,
		port)
	log.Fatal(srv.ListenAndServe())
}
