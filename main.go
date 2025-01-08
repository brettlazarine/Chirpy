package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %v on port: %v", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
	defer srv.Close()
}
