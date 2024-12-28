### Task Management System

## Overview

The **Task Management System** is a Golang-based application that allows users to manage tasks within **workspaces**. Each user can:

- Create up to 6 **workspaces**.
- Be invited to 6 other **workspaces**.
- Assign and manage tasks with other users.
- Track task status and updates.

This project is designed to streamline task management, improve collaboration, and offer a user-friendly experience for team-based projects.

## Features

- **Workspaces**:
    - Users can create and manage up to 6 personal workspaces.
    - Collaboration is enabled by inviting other users to join a workspace.

- **Tasks**:
    - Assign tasks to specific users in a **workspace**.
    - Manage task status: To Do, In Progress, or Completed.
    - View task history and updates.
- **User Roles**:
    - Admin: Full control over a **workspace** and **tasks**.
    - Member: Can view, create, and update tasks within a **workspace**.
- **Security**:
    - Ownership validation to prevent unauthorized access or modifications.
    - CSRF protection and session management.


## Installation
### Prerequisites
- Go (version 1.22.1)
- Docker and Docker Compose
- A terminal with zsh or bash

### Steps
1. Clone the repository:

```bash
# With SSH
git clone git@github.com:andres085/golang-task-manager.git

# With HTTPS
git clone https://github.com/andres085/golang-task-manager.git
``` 

2. Start docker container
```bash
cd task-manager

docker compose up -d
```

3. Run the migrations to set up the database:
```bash
go run ./cmd/migrate
```

4. Start the application:
```bash
go run ./cmd/web
```

## Usage

**Register**
- Enter the register link, fill out the form, and you will be redirected to the login page.

**Login**
- Log in with the credentials created earlier. If the login is successful, you will be redirected to the workspaces view.

**Workspaces**
- **Create a workspace**: Click on the "Create Workspace" button.
- **Invite users**: Use the "Add Users" button in the workspace view (redirected after creating a new one or clicking on the workspace title). Search for a user by email and click on "Add User" to invite them.
- **Remove a user**: In the workspace view go to "Add User" click "Remove". The admin can't be removed.
- **Update a workspace**: Workspace owners can update a workspace after validating ownership.
- **Delete a workspace**: Workspace owners can delete a workspace after validating ownership.
- **Invited Workspaces**: Go to "Invited Workspaces" tab to see the workspaces where you have been invited.

**Tasks**
- **Create a task**: Navigate to the workspace view or the workspace detail view, select "View Tasks" or "Add Task" to go to the tasks view to have access to the "Create Task" button.
- **View task list**: Access the tasks view page to see the list of tasks by clicking in workspaces "View Tasks".
- **Task Detail**: Click on the task title to go to the task view and see the task details.
- **Update a task**: Modify task data, status, or reassign the task by clicking on "Edit" in the task view page or the table.
- **Delete a task**: Delete a task by clicking on the "Delete" button on the table or in the task view.

## Testing
The project has a combination of unit tests and integration tests using a dockerized test database.

Run tests using:
```bash
go test ./...
```
Coverage of handlers around 80%, open `coverage.html` file to check it out.

## Future Enhancements
- Add the possibility to comment on tasks inside the task view.
- Implement notifications for task updates and deadlines.
- Improve styles.
- Add user profiles.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments
- Inspired by Trello for task management concepts.
- Special thanks to me, the only contributor of this project, and my cats for the support.
