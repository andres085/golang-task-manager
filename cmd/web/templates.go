package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/andres085/task_manager/internal/models"
	"github.com/andres085/task_manager/ui"
)

type templateData struct {
	CurrentYear     int
	Task            models.Task
	Tasks           []models.Task
	Workspace       models.Workspace
	Workspaces      []models.Workspace
	User            *models.User
	WorkspaceUsers  []models.UserWithRole
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/nav.html",
			"html/partials/confirmation_modal.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
