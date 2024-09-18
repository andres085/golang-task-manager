package main

import (
	"net/http"

	"github.com/andres085/task_manager/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", app.ping)

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /task/view/{id}", app.taskView)
	mux.HandleFunc("GET /task/update/{id}", app.taskUpdate)
	mux.HandleFunc("POST /task/create", app.taskCreatePost)
	mux.HandleFunc("POST /task/update/{id}", app.taskUpdatePost)
	mux.HandleFunc("POST /task/delete/{id}", app.taskDelete)

	mux.HandleFunc("GET /workspace/view", app.workspaceViewAll)
	mux.HandleFunc("GET /workspace/view/{id}", app.workspaceView)
	mux.HandleFunc("GET /workspace/view/{id}/tasks", app.taskViewAll)
	mux.HandleFunc("GET /workspace/create", app.workspaceCreate)
	mux.HandleFunc("GET /workspace/update/{id}", app.workspaceUpdate)
	mux.HandleFunc("GET /workspace/{id}/task/create", app.taskCreate)
	mux.HandleFunc("POST /workspace/create", app.workspaceCreatePost)
	mux.HandleFunc("POST /workspace/update/{id}", app.workspaceUpdatePost)
	mux.HandleFunc("POST /workspace/delete/{id}", app.workspaceDelete)

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
