package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/andres085/task_manager/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

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

	data := app.newTemplateData(r)
	data.Task = task

	app.render(w, r, http.StatusOK, "task_view.html", data)
}

func (app *application) taskViewAll(w http.ResponseWriter, r *http.Request) {

	tasks, err := app.tasks.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Tasks = tasks

	app.render(w, r, http.StatusOK, "tasks_view.html", data)
}

type taskCreateForm struct {
	Title       string
	Content     string
	Priority    string
	FieldErrors map[string]string
}

func (app *application) taskCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = taskCreateForm{
		Priority: "LOW",
	}

	app.render(w, r, http.StatusOK, "task_create.html", data)
}

func (app *application) taskCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := taskCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Priority:    r.PostForm.Get("priority"),
		FieldErrors: make(map[string]string),
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "task_create.html", data)
		return
	}

	id, err := app.tasks.Insert(form.Title, form.Content, form.Priority)
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

	http.Redirect(w, r, "/task/view", http.StatusSeeOther)
}
