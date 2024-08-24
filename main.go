package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the task manager"))
}

func taskView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display a specific task with ID %d...", id)
	w.Write([]byte(msg))
}

func taskCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display form to create a new task..."))
}

func taskCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Save a new task..."))
}

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
