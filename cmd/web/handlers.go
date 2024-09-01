package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andres085/task_manager/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	app.render(w, r, http.StatusOK, "home.html", templateData{})
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

	data := templateData{
		Task: task,
	}

	app.render(w, r, http.StatusOK, "task_view.html", data)
}

func (app *application) taskViewAll(w http.ResponseWriter, r *http.Request) {

	tasks, err := app.tasks.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Tasks: tasks,
	}

	app.render(w, r, http.StatusOK, "tasks_view.html", data)
}

func (app *application) taskCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new task..."))
}

func (app *application) taskCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "Task from the backend"
	content := "Task from the backend content test asd 123"
	priority := "LOW"

	id, err := app.tasks.Insert(title, content, priority)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/task/view/%d", id), http.StatusSeeOther)
}
