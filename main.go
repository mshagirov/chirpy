package main

import (
	"log"
	"net/http"
)

func main() {
	const rootPath = "."
	const port = "8080"
	apiCfg := apiConfig{}
	apiCfg.fileserverHits.Store(0)
	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
	mux.Handle("/app/", fsHandler)
	mux.HandleFunc("GET /api/healthz", ServerReady)
	mux.HandleFunc("GET /admin/metrics", apiCfg.serveMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.serveReset)
	mux.HandleFunc("POST /api/validate_chirp", serveValidateChirp)
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s port: %s\n",
		rootPath,
		port)
	log.Fatal(srv.ListenAndServe())
}
