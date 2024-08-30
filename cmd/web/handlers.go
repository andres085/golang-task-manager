package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/andres085/task_manager/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}
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

	fmt.Fprintf(w, "%+v", task)
}

func (app *application) taskViewAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := app.tasks.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, task := range tasks {
		fmt.Fprintf(w, "%+v\n", task)
	}
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
