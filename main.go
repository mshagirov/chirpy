package main

import (
	// "io"
	"log"
	"net/http"
)

func main() {
	const rootPath = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(rootPath))))
	mux.HandleFunc("/healthz", ServerReady)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s port: %s\n",
		rootPath,
		port)
	log.Fatal(srv.ListenAndServe())
}

func ServerReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	n, err := w.Write([]byte("200 " + http.StatusText(http.StatusOK)))
	if err != nil {
		log.Printf("Error: responseWriter returned %v, %v", n, err)
	}
}
