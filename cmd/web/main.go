package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /task/view/{id}", taskView)
	mux.HandleFunc("GET /task/create", taskCreate)
	mux.HandleFunc("POST /task/create", taskCreatePost)

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
