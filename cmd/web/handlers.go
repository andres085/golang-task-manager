package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andres085/task_manager/internal/models"
	"github.com/andres085/task_manager/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) taskView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	task, err := app.tasks.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Task = task

	data.Flash = flash

	app.render(w, r, http.StatusOK, "task_view.html", data)
}

func (app *application) taskViewAll(w http.ResponseWriter, r *http.Request) {
	workspaceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	tasks, err := app.tasks.GetAll(workspaceId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Tasks = tasks
	data.Workspace.ID = workspaceId

	app.render(w, r, http.StatusOK, "tasks_view.html", data)
}

type taskCreateForm struct {
	ID                  *int
	Title               string `form:"title"`
	Content             string `form:"content"`
	Priority            string `form:"priority"`
	WorkspaceID         int    `form:"workspace_id"`
	validator.Validator `form:"-"`
}

func (app *application) taskCreate(w http.ResponseWriter, r *http.Request) {
	workspaceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	data := app.newTemplateData(r)

	data.Form = taskCreateForm{
		Priority:    "LOW",
		WorkspaceID: workspaceId,
	}

	app.render(w, r, http.StatusOK, "task_create.html", data)
}

func (app *application) taskCreatePost(w http.ResponseWriter, r *http.Request) {
	var form taskCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "task_create.html", data)
		return
	}

	id, err := app.tasks.Insert(form.Title, form.Content, form.Priority, form.WorkspaceID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Task successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/task/view/%d", id), http.StatusSeeOther)
}

func (app *application) taskUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	task, err := app.tasks.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)

	data.Form = taskCreateForm{
		ID:       &task.ID,
		Title:    task.Title,
		Content:  task.Content,
		Priority: task.Priority,
	}

	app.render(w, r, http.StatusOK, "task_update.html", data)
}

func (app *application) taskUpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	var form taskCreateForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		form.ID = &id
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "task_update.html", data)
		return
	}

	err = app.tasks.Update(id, form.Title, form.Content, form.Priority)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/task/view/%d", id), http.StatusSeeOther)
}

func (app *application) taskDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspaceId, err := strconv.Atoi(r.PathValue("workspaceId"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	row, err := app.tasks.Delete(id)
	if err != nil || row < 1 {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workspace/view/%d/tasks", workspaceId), http.StatusSeeOther)
}

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

type workspaceCreateForm struct {
	ID                  *int
	Title               string `form:"title"`
	Description         string `form:"description"`
	validator.Validator `form:"-"`
}

func (app *application) workspaceCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = workspaceCreateForm{}

	app.render(w, r, http.StatusOK, "workspace_create.html", data)
}

func (app *application) workspaceCreatePost(w http.ResponseWriter, r *http.Request) {
	var form workspaceCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Description, 255), "description", "This field cannot be more than 255 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "workspace_create.html", data)
		return
	}

	id, err := app.workspaces.Insert(form.Title, form.Description)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Workspace successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/workspace/view/%d", id), http.StatusSeeOther)
}

func (app *application) workspaceView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	workspace, err := app.workspaces.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData(r)
	data.Workspace = workspace
	data.Flash = flash

	app.render(w, r, http.StatusOK, "workspace_view.html", data)
}

func (app *application) workspaceViewAll(w http.ResponseWriter, r *http.Request) {
	workspaces, err := app.workspaces.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Workspaces = workspaces

	app.render(w, r, http.StatusOK, "workspaces_view.html", data)
}

func (app *application) workspaceUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	workspace, err := app.workspaces.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)

	data.Form = workspaceCreateForm{
		ID:          &workspace.ID,
		Title:       workspace.Title,
		Description: workspace.Description,
	}

	app.render(w, r, http.StatusOK, "workspace_update.html", data)
}

func (app *application) workspaceUpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	var form workspaceCreateForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Description), "description", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Description, 255), "description", "This field cannot be more than 255 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		form.ID = &id
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "workspace_update.html", data)
		return
	}

	err = app.workspaces.Update(id, form.Title, form.Description)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workspace/view/%d", id), http.StatusSeeOther)
}

func (app *application) workspaceDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	row, err := app.workspaces.Delete(id)
	if err != nil || row < 1 {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "/workspace/view", http.StatusSeeOther)
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignIn(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "signin.html", data)
}
