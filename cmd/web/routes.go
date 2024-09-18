package main

import (
	"net/http"

	"github.com/andres085/task_manager/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /ping", dynamic.ThenFunc(app.ping))
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /task/view/{id}", dynamic.ThenFunc(app.taskView))
	mux.Handle("GET /task/update/{id}", dynamic.ThenFunc(app.taskUpdate))
	mux.Handle("POST /task/create", dynamic.ThenFunc(app.taskCreatePost))
	mux.Handle("POST /task/update/{id}", dynamic.ThenFunc(app.taskUpdatePost))
	mux.Handle("POST /task/delete/{id}", dynamic.ThenFunc(app.taskDelete))

	mux.Handle("GET /workspace/view", dynamic.ThenFunc(app.workspaceViewAll))
	mux.Handle("GET /workspace/view/{id}", dynamic.ThenFunc(app.workspaceView))
	mux.Handle("GET /workspace/view/{id}/tasks", dynamic.ThenFunc(app.taskViewAll))
	mux.Handle("GET /workspace/create", dynamic.ThenFunc(app.workspaceCreate))
	mux.Handle("GET /workspace/update/{id}", dynamic.ThenFunc(app.workspaceUpdate))
	mux.Handle("GET /workspace/{id}/task/create", dynamic.ThenFunc(app.taskCreate))
	mux.Handle("POST /workspace/create", dynamic.ThenFunc(app.workspaceCreatePost))
	mux.Handle("POST /workspace/update/{id}", dynamic.ThenFunc(app.workspaceUpdatePost))
	mux.Handle("POST /workspace/delete/{id}", dynamic.ThenFunc(app.workspaceDelete))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
