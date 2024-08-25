package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Write([]byte("Hello from the Taskmanager"))
}

func taskView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific task with ID %d...", id)
}

func taskCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new task..."))
}
func taskCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new task..."))
}
