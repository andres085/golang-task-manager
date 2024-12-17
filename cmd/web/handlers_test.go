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

	t.Run("Unauthenticated", func(t *testing.T) {

		code, headers, _ := ts.get(t, "/workspace/view/1/tasks")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	t.Run("NotFound", func(t *testing.T) {
		ts.loginUser(t)
		code, _, _ := ts.get(t, "/workspace/view/5/tasks")

		assert.Equal(t, code, http.StatusNotFound)
	})

	t.Run("Authenticated", func(t *testing.T) {
		ts.loginUser(t)

		code, _, body := ts.get(t, "/workspace/view/1/tasks")
		wantTitle := "Tasks View"
		firstTestTaskTitle := "First Test Task"
		secondTestTaskTitle := "Second Test Task"

		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, wantTitle)
		assert.StringContains(t, body, firstTestTaskTitle)
		assert.StringContains(t, body, secondTestTaskTitle)
	})
}

func TestTaskView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/task/view/1")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	ts.loginUser(t)

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

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/workspace/1/task/create")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		ts.loginUser(t)

		code, _, body := ts.get(t, "/workspace/1/task/create")
		wantTitle := "Create New Task"

		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, wantTitle)
	})
}

func TestTaskCreatePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name      string
		title     string
		content   string
		priority  string
		csrfToken string
		wantCode  int
	}{
		{
			name:      "Valid Submission",
			title:     "Test Task",
			content:   "Test Content",
			priority:  "LOW",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:      "Invalid Submission without Title",
			title:     "",
			content:   "Test Content",
			priority:  "LOW",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid Submission without Content",
			title:     "Test Task",
			content:   "",
			priority:  "LOW",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("content", tt.content)
			form.Add("priority", tt.priority)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/task/create", form)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestTaskUpdate(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/task/update/1")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	t.Run("NotFound", func(t *testing.T) {
		ts.loginUser(t)
		code, _, _ := ts.get(t, "/task/update/99")

		assert.Equal(t, code, http.StatusNotFound)
	})

	t.Run("NotFound negative id", func(t *testing.T) {
		ts.loginUser(t)
		code, _, _ := ts.get(t, "/task/update/-1")

		assert.Equal(t, code, http.StatusNotFound)
	})

	t.Run("Authenticated", func(t *testing.T) {
		ts.loginUser(t)

		code, _, body := ts.get(t, "/task/update/1")
		titleInput := `<input type="text" class="form-control " id="title" name="title"
        placeholder="Enter task title" value="First Test Task">`

		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, titleInput)
	})
}

func TestTaskUpdatePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name      string
		title     string
		content   string
		priority  string
		csrfToken string
		wantCode  int
		urlPath   string
	}{
		{
			name:      "Valid Submission",
			title:     "Test Task",
			content:   "Test Content",
			priority:  "LOW",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
			urlPath:   "/task/update/1",
		},
		{
			name:      "Invalid Submission without Title",
			title:     "",
			content:   "Test Content",
			priority:  "LOW",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
			urlPath:   "/task/update/1",
		},
		{
			name:      "Invalid Submission without Content",
			title:     "Test Task",
			content:   "",
			priority:  "LOW",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
			urlPath:   "/task/update/1",
		},
		{
			name:      "Invalid Task Id",
			title:     "Test Task",
			content:   "Test Content",
			priority:  "LOW",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusNotFound,
			urlPath:   "/task/update/-1",
		},
		{
			name:      "Invalid task owner",
			title:     "Test Task",
			content:   "Test Content",
			priority:  "LOW",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusNotFound,
			urlPath:   "/task/update/2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("content", tt.content)
			form.Add("priority", tt.priority)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, tt.urlPath, form)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestTaskDelete(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", validCSRFToken)

	code, _, _ := ts.get(t, "/workspace/1/task/delete/1")

	assert.Equal(t, code, http.StatusMethodNotAllowed)
}

func TestTaskDeletePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", validCSRFToken)

	code, _, _ := ts.postForm(t, "/workspace/1/task/delete/1", form)

	assert.Equal(t, code, http.StatusSeeOther)
}

func TestWorkspaceViewAll(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	code, _, body := ts.get(t, "/workspace/view")
	wantTitle := "Workspaces View"
	firstTestTaskTitle := "First Workspace"
	secondTestTaskTitle := "Second Workspace"

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, wantTitle)
	assert.StringContains(t, body, firstTestTaskTitle)
	assert.StringContains(t, body, secondTestTaskTitle)
}

func TestWorkspaceView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/workspace/view/1")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	ts.loginUser(t)

	tests := []struct {
		name            string
		urlPath         string
		wantCode        int
		wantTitle       string
		wantDescription string
	}{
		{
			name:            "Valid ID",
			urlPath:         "/workspace/view/1",
			wantCode:        http.StatusOK,
			wantTitle:       "First Workspace",
			wantDescription: "First workspace Description",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/workspace/view/3",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/workspace/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/workspace/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/workspace/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/workspace/view/",
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

			if tt.wantDescription != "" {
				assert.StringContains(t, body, tt.wantDescription)
			}
		})
	}
}

func TestWorkspaceCreate(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	code, _, body := ts.get(t, "/workspace/create")
	titleInput := `<input type="text" class="form-control " id="title" name="title"
        placeholder="Enter workspace title" value="">`

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, titleInput)
}

func TestWorkspaceCreatePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		title       string
		description string
		csrfToken   string
		wantCode    int
	}{
		{
			name:        "Valid Submission",
			title:       "Test Workspace",
			description: "Test workspace description",
			csrfToken:   validCSRFToken,
			wantCode:    http.StatusSeeOther,
		},
		{
			name:        "Invalid Submission without Title",
			title:       "",
			description: "Test workspace description",
			csrfToken:   validCSRFToken,
			wantCode:    http.StatusUnprocessableEntity,
		},
		{
			name:        "Invalid Submission without description",
			title:       "Test Task",
			description: "",
			csrfToken:   validCSRFToken,
			wantCode:    http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("description", tt.description)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/workspace/create", form)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestWorkspaceUpdate(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	code, _, body := ts.get(t, "/workspace/update/1")
	titleInput := `<input type="text" class="form-control " id="title" name="title"
        placeholder="Enter workspace title" value="First Workspace">`

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, titleInput)
}

func TestWorkspaceUpdatePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	ts.loginUser(t)

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name        string
		title       string
		description string
		csrfToken   string
		wantCode    int
	}{
		{
			name:        "Valid Submission",
			title:       "Test Workspace",
			description: "Test workspace description",
			csrfToken:   validCSRFToken,
			wantCode:    http.StatusSeeOther,
		},
		{
			name:        "Invalid Submission without Title",
			title:       "",
			description: "Test workspace description",
			csrfToken:   validCSRFToken,
			wantCode:    http.StatusUnprocessableEntity,
		},
		{
			name:        "Invalid Submission without description",
			title:       "Test Task",
			description: "",
			csrfToken:   validCSRFToken,
			wantCode:    http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("description", tt.description)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/workspace/update/1", form)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestWorkspaceDeletePost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/user/login")
		validCSRFToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("csrf_token", validCSRFToken)

		code, headers, _ := ts.postForm(t, "/workspace/delete/1", form)

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		ts.loginUser(t)

		_, _, body := ts.get(t, "/user/login")
		validCSRFToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("csrf_token", validCSRFToken)

		code, _, _ := ts.postForm(t, "/workspace/delete/1", form)

		assert.Equal(t, code, http.StatusSeeOther)
	})

}

func TestWorkspaceAddUserView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/workspace/1/user/add")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	ts.loginUser(t)

	tests := []struct {
		name      string
		urlPath   string
		wantCode  int
		wantTitle string
	}{
		{
			name:      "Valid ID",
			urlPath:   "/workspace/1/user/add",
			wantCode:  http.StatusOK,
			wantTitle: "Users in this Workspace",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/workspace/5/user/add",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/workspace/-1/user/add",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/workspace/1.23/user/add",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/workspace/a/user/add",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/workspace/user/add",
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
		})
	}
}

func TestUserRegisterHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/user/register")
	formTitle := `<h2 class="mb-4 text-center">Register</h2>`

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, formTitle)
}

func TestUserRegisterPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name      string
		firstName string
		lastName  string
		email     string
		password  string
		csrfToken string
		wantCode  int
	}{
		{
			name:      "Valid Submission",
			firstName: "Test",
			lastName:  "McTester",
			email:     "test@mail.com",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:      "Invalid submission without Username",
			firstName: "",
			lastName:  "McTester",
			email:     "test@mail.com",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid Username",
			firstName: "TestTestTestTestTestTest",
			lastName:  "McTester",
			email:     "test@mail.com",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid Lastname",
			firstName: "Test",
			lastName:  "McTesterMcTesterMcTesterMcTesterMcTesterMcTester",
			email:     "test@mail.com",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid submission without LastName",
			firstName: "TestTestTestTestTestTest",
			lastName:  "",
			email:     "test@mail.com",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid Email",
			firstName: "Test",
			lastName:  "McTester",
			email:     "testmail.com",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid submission without Email",
			firstName: "Test",
			lastName:  "McTester",
			email:     "",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid Password",
			firstName: "Test",
			lastName:  "McTester",
			email:     "test@mail.com",
			password:  "pa$$",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid submission without Password",
			firstName: "Test",
			lastName:  "McTester",
			email:     "test@mail.com",
			password:  "",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("firstName", tt.firstName)
			form.Add("lastName", tt.lastName)
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/user/register", form)

			assert.Equal(t, code, tt.wantCode)
		})
	}

}

func TestUserLoginHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/user/login")
	formTitle := `<h2 class="mb-4 text-center">Login</h2>`

	assert.Equal(t, code, http.StatusOK)
	assert.StringContains(t, body, formTitle)
}

func TestUserLoginPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	tests := []struct {
		name      string
		email     string
		password  string
		csrfToken string
		wantCode  int
	}{
		{
			name:      "Valid Submission",
			email:     "alice@example.com",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:      "Invalid Email",
			email:     "testmail.com",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid submission without Email",
			email:     "",
			password:  "pa$$word",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
		{
			name:      "Invalid submission without Password",
			email:     "test@mail.com",
			password:  "",
			csrfToken: validCSRFToken,
			wantCode:  http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/user/login", form)

			assert.Equal(t, code, tt.wantCode)
		})
	}

}

func TestUserLogoutPost(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/login")
	validCSRFToken := extractCSRFToken(t, body)

	form := url.Values{}
	form.Add("csrf_token", validCSRFToken)

	code, _, _ := ts.postForm(t, "/user/logout", form)

	assert.Equal(t, code, http.StatusSeeOther)
}
