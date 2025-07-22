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
	mux.Handle("/app/",
		http.StripPrefix("/app",
			apiCfg.middlewareMetricsInc(
				http.FileServer(http.Dir(rootPath)),
			),
		),
	)
	mux.HandleFunc("/healthz", ServerReady)
	mux.HandleFunc("/metrics", apiCfg.serveMetrics)
	mux.HandleFunc("/reset", apiCfg.serveReset)
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s port: %s\n",
		rootPath,
		port)
	log.Fatal(srv.ListenAndServe())
}
