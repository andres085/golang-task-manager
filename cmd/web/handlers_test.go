package main

import (
	"net/http"
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
