package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strconv"
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
	Limit           int
	CurrentPage     int
	TotalPages      int
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func iterPages(total int) []int {
	var pages []int
	for i := 1; i <= total; i++ {
		pages = append(pages, i)
	}
	return pages
}

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"iterPages": iterPages,
	"add":       add,
	"sub":       sub,
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

func getPaginationParams(r *http.Request, defaultLimit int) (int, int, int) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * limit

	return limit, page, offset
}
