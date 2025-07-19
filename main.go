package main

import (
	"log"
	"net/http"
)

func main() {
	const rootPath = "."
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(rootPath)))
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("Serving files from %s port: %s\n",
		rootPath,
		port)
	log.Fatal(srv.ListenAndServe())
}
