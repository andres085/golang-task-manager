package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/andres085/task_manager/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}

func TestHomeView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/")
	wantTitle := "Welcome to the Task Manager"

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, wantTitle)
}

func TestTaskViewAll(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/workspace/view/1/tasks")
	wantTitle := "Tasks View"
	firstTestTaskTitle := "First Test Task"
	secondTestTaskTitle := "Second Test Task"

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, wantTitle)
	assert.StringContains(t, body, firstTestTaskTitle)
	assert.StringContains(t, body, secondTestTaskTitle)
}

func TestTaskView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name        string
		urlPath     string
		wantCode    int
		wantTitle   string
		wantContent string
	}{
		{
			name:        "Valid ID",
			urlPath:     "/task/view/1",
			wantCode:    http.StatusOK,
			wantTitle:   "First Test Task",
			wantContent: "First Test Task Content",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/task/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/task/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/task/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/task/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/task/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantTitle != "" {
				assert.StringContains(t, body, tt.wantTitle)
			}

			if tt.wantContent != "" {
				assert.StringContains(t, body, tt.wantContent)
			}
		})
	}
}

func TestTaskCreate(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/workspace/1/task/create")
	wantTitle := "Create New Task"

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, wantTitle)
}

func TestTaskCreatePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		title    string
		content  string
		priority string
		wantCode int
	}{
		{
			name:     "Valid Submission",
			title:    "Test Task",
			content:  "Test Content",
			priority: "LOW",
			wantCode: http.StatusSeeOther,
		},
		{
			name:     "Invalid Submission without Title",
			title:    "",
			content:  "Test Content",
			priority: "LOW",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Invalid Submission without Content",
			title:    "Test Task",
			content:  "",
			priority: "LOW",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("content", tt.content)
			form.Add("priority", tt.priority)

			code, _, _ := ts.postForm(t, "/task/create", form)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestTaskUpdate(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/task/update/1")
	titleInput := `<input type="text" class="form-control " id="title" name="title"
        placeholder="Enter task title" value="First Test Task">`

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, titleInput)
}

func TestTaskUpdatePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		title    string
		content  string
		priority string
		wantCode int
	}{
		{
			name:     "Valid Submission",
			title:    "Test Task",
			content:  "Test Content",
			priority: "LOW",
			wantCode: http.StatusSeeOther,
		},
		{
			name:     "Invalid Submission without Title",
			title:    "",
			content:  "Test Content",
			priority: "LOW",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "Invalid Submission without Content",
			title:    "Test Task",
			content:  "",
			priority: "LOW",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("content", tt.content)
			form.Add("priority", tt.priority)

			code, _, _ := ts.postForm(t, "/task/update/1", form)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestTaskDeletePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	form := url.Values{}

	code, _, _ := ts.postForm(t, "/task/delete/1", form)

	assert.Equal(t, code, http.StatusSeeOther)
}

func TestWorkspaceViewAll(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/workspace/view")
	wantTitle := "Workspaces View"
	firstTestTaskTitle := "First Workspace"
	secondTestTaskTitle := "Second Workspace"

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, wantTitle)
	assert.StringContains(t, body, firstTestTaskTitle)
	assert.StringContains(t, body, secondTestTaskTitle)
}
