package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /ping", app.ping)

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /task/view", app.taskViewAll)
	mux.HandleFunc("GET /task/view/{id}", app.taskView)
	mux.HandleFunc("GET /task/create", app.taskCreate)
	mux.HandleFunc("GET /task/update/{id}", app.taskUpdate)
	mux.HandleFunc("POST /task/create", app.taskCreatePost)
	mux.HandleFunc("POST /task/delete/{id}", app.taskDelete)
	mux.HandleFunc("POST /task/update/{id}", app.taskUpdatePost)

	return mux
}
