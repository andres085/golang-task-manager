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

	data := app.newTemplateData(r)
	data.Task = task

	app.render(w, r, http.StatusOK, "task_view.html", data)
}

func (app *application) taskViewAll(w http.ResponseWriter, r *http.Request) {
	workspaceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	tasks, err := app.tasks.GetAll(workspaceId, limit, offset)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Tasks = tasks
	data.Workspace.ID = workspaceId
	data.Limit = limit

	app.render(w, r, http.StatusOK, "tasks_view.html", data)
}

type taskCreateForm struct {
	ID                  *int
	Title               string                `form:"title"`
	Content             string                `form:"content"`
	Priority            string                `form:"priority"`
	WorkspaceID         int                   `form:"workspace_id"`
	UserID              int                   `form:"user_id"`
	DefaultUser         models.UserWithRole   `form:"-"`
	WorkspaceUsers      []models.UserWithRole `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) taskCreate(w http.ResponseWriter, r *http.Request) {
	workspaceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	data := app.newTemplateData(r)

	workspaceUsers, err := app.users.GetWorkspaceUsers(workspaceId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var adminUser models.UserWithRole
	regularUsers := make([]models.UserWithRole, len(workspaceUsers)-1)

	for i, user := range workspaceUsers {
		if user.Role == "ADMIN" {
			adminUser = user
			regularUsers = append(workspaceUsers[:i], workspaceUsers[i+1:]...)
		}
	}

	data.Form = taskCreateForm{
		Priority:       "LOW",
		WorkspaceID:    workspaceId,
		DefaultUser:    adminUser,
		WorkspaceUsers: regularUsers,
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

	id, err := app.tasks.Insert(form.Title, form.Content, form.Priority, form.WorkspaceID, form.UserID)
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

	workspaceUsers, err := app.users.GetWorkspaceUsers(task.WorkspaceId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var assignedUser models.UserWithRole
	otherUsers := make([]models.UserWithRole, len(workspaceUsers)-1)

	for i, user := range workspaceUsers {
		if user.ID == task.UserId {
			assignedUser = user
			otherUsers = append(workspaceUsers[:i], workspaceUsers[i+1:]...)
		}
	}

	data := app.newTemplateData(r)

	data.Form = taskCreateForm{
		ID:             &task.ID,
		Title:          task.Title,
		Content:        task.Content,
		Priority:       task.Priority,
		DefaultUser:    assignedUser,
		WorkspaceUsers: otherUsers,
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

	err = app.tasks.Update(id, form.Title, form.Content, form.Priority, form.UserID)
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

	userId := r.Context().Value(userIDContextKey).(int)

	id, err := app.workspaces.Insert(form.Title, form.Description, userId)
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

	data := app.newTemplateData(r)
	data.Workspace = workspace

	app.render(w, r, http.StatusOK, "workspace_view.html", data)
}

func (app *application) workspaceViewAll(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userIDContextKey).(int)
	workspaces, err := app.workspaces.GetAll(userId)
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
	userId, ok := r.Context().Value(userIDContextKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workspaceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || workspaceId < 1 {
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
		form.ID = &workspaceId
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "workspace_update.html", data)
		return
	}

	isOwner, err := app.workspaces.ValidateOwnership(userId, workspaceId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !isOwner {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = app.workspaces.Update(workspaceId, form.Title, form.Description)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workspace/view/%d", workspaceId), http.StatusSeeOther)
}

func (app *application) workspaceAddUser(w http.ResponseWriter, r *http.Request) {
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

	email := r.URL.Query().Get("email")
	var foundUser *models.User

	if email != "" {
		foundUser, err = app.users.GetUserToInvite(email, workspace.ID)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.sessionManager.Put(r.Context(), "flash", "User not found or already added")
				http.Redirect(w, r, fmt.Sprintf("/workspace/%d/user/add", workspace.ID), http.StatusSeeOther)
			} else {
				app.serverError(w, r, err)
			}
			return
		}
	}

	workspaceUsers, err := app.users.GetWorkspaceUsers(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Workspace = workspace
	data.User = foundUser
	data.WorkspaceUsers = workspaceUsers

	app.render(w, r, http.StatusOK, "workspace_users.html", data)
}

type addUserForm struct {
	UserID int `form:"userID"`
}

func (app *application) workspaceAddUserPost(w http.ResponseWriter, r *http.Request) {
	workspaceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	var form addUserForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(userIDContextKey).(int)
	if userId == form.UserID {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = app.users.AddUserToWorkspace(form.UserID, workspaceId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workspace/%d/user/add", workspaceId), http.StatusSeeOther)
}

func (app *application) workspaceRemoveUserPost(w http.ResponseWriter, r *http.Request) {
	workspaceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	userId, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	row, err := app.users.RemoveUserFromWorkspace(workspaceId, userId)
	if err != nil || row < 1 {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/workspace/%d/user/add", workspaceId), http.StatusSeeOther)
}

func (app *application) workspaceDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	workspaceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || workspaceId < 1 {
		http.NotFound(w, r)
		return
	}

	row, err := app.workspaces.Delete(workspaceId)
	if err != nil || row < 1 {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "/workspace/view", http.StatusSeeOther)
}

type userCreateForm struct {
	FirstName           string `form:"firstName"`
	LastName            string `form:"lastName"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = userCreateForm{}

	app.render(w, r, http.StatusOK, "register.html", data)
}

func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	var form userCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.FirstName), "firstName", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.FirstName, 20), "firstName", "This field cannot be more than 20 characters long")
	form.CheckField(validator.NotBlank(form.LastName), "lastName", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.LastName, 20), "lastName", "This field cannot be more than 20 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 6), "password", "This field cannot be less than 6 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "register.html", data)
		return
	}

	err = app.users.Insert(form.FirstName, form.LastName, form.Email, form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User registered successfully!")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/workspace/view", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
