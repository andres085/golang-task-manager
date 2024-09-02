package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /task/view", app.taskViewAll)
	mux.HandleFunc("GET /task/view/{id}", app.taskView)
	mux.HandleFunc("GET /task/create", app.taskCreate)
	mux.HandleFunc("POST /task/create", app.taskCreatePost)
	mux.HandleFunc("POST /task/delete/{id}", app.taskDelete)

	return mux
}
