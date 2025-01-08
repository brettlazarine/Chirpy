package main

import (
	"net/http"
	"log"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
	defer srv.Close()
}