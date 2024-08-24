package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the task manager"))
}

func taskView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific task..."))
}

func taskCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display form to create a new task..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/task/view", taskView)
	mux.HandleFunc("/task/create", taskCreate)

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
