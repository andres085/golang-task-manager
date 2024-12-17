package main

import (
	"net/http"

	"github.com/andres085/task_manager/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	protected := dynamic.Append(app.requireAuthentication)
	workspaceMembership := protected.Append(app.checkWorkspaceMembership)
	workspaceAdminPermission := protected.Append(app.checkWorkspaceAdmin)
	taskAdminPermission := protected.Append(app.checkTaskAdmin)

	mux.HandleFunc("GET /ping", app.ping)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))

	mux.Handle("GET /task/view/{id}", protected.ThenFunc(app.taskView))
	mux.Handle("GET /task/update/{id}", protected.ThenFunc(app.taskUpdate))
	mux.Handle("GET /workspace/{id}/task/create", workspaceMembership.ThenFunc(app.taskCreate))
	mux.Handle("POST /task/create", protected.ThenFunc(app.taskCreatePost))
	mux.Handle("POST /task/update/{id}", protected.ThenFunc(app.taskUpdatePost))
	mux.Handle("POST /workspace/{workspaceId}/task/delete/{id}", taskAdminPermission.ThenFunc(app.taskDelete))

	mux.Handle("GET /workspace/view", protected.ThenFunc(app.workspaceViewAll))
	mux.Handle("GET /workspace/view/{id}", workspaceMembership.ThenFunc(app.workspaceView))
	mux.Handle("GET /workspace/view/{id}/tasks", protected.ThenFunc(app.taskViewAll))
	mux.Handle("GET /workspace/create", protected.ThenFunc(app.workspaceCreate))
	mux.Handle("GET /workspace/update/{id}", workspaceAdminPermission.ThenFunc(app.workspaceUpdate))
	mux.Handle("GET /workspace/{id}/user/add", workspaceAdminPermission.ThenFunc(app.workspaceAddUser))
	mux.Handle("POST /workspace/create", protected.ThenFunc(app.workspaceCreatePost))
	mux.Handle("POST /workspace/update/{id}", protected.ThenFunc(app.workspaceUpdatePost))
	mux.Handle("POST /workspace/delete/{id}", workspaceAdminPermission.ThenFunc(app.workspaceDelete))
	mux.Handle("POST /workspace/{id}/user/add", workspaceMembership.ThenFunc(app.workspaceAddUserPost))
	mux.Handle("POST /workspace/{id}/user/remove/{userId}", workspaceMembership.ThenFunc(app.workspaceRemoveUserPost))

	mux.Handle("GET /user/register", dynamic.ThenFunc(app.userSignUp))
	mux.Handle("POST /user/register", dynamic.ThenFunc(app.userSignUpPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
