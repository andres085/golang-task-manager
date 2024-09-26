CREATE TABLE users (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	firstName VARCHAR(255) NOT NULL,
	lastName VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

CREATE TABLE tasks (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    priority TEXT NOT NULL,
    created DATETIME NOT NULL,
    finished DATETIME NOT NULL,
    workspace_id INTEGER NOT NULL,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
);

CREATE INDEX idx_tasks_created ON tasks(created);


CREATE TABLE workspaces (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL UNIQUE,
    description TEXT(255) NOT NULL,
    created DATETIME NOT NULL
);

CREATE INDEX idx_workspace_id ON workspaces(id);

CREATE TABLE users_workspaces (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	user_id INT NOT NULL,
	workspace_id INT NOT NULL,
	role TEXT NOT NULL,
	created DATETIME NOT NULL,
	UNIQUE(user_id, workspace_id),
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
)

CREATE INDEX idx_user_id ON users_workspaces (user_id);
CREATE INDEX idx_workspace_id ON users_workspaces (workspace_id);

INSERT INTO workspaces (title, description, created) VALUES (
    'First Workspace',
    'This is the first workspace description',
    UTC_TIMESTAMP()
);

INSERT INTO workspaces (title, description, created) VALUES (
    'Second Workspace',
    'This is the second workspace description',
    UTC_TIMESTAMP()
);

INSERT INTO workspaces (title, description, created) VALUES (
    'Third Workspace',
    'This is the third workspace description',
    UTC_TIMESTAMP()
);

INSERT INTO tasks (title, content, priority, created, finished, workspace_id) VALUES (
    'First Task',
    'This is the content of the first task',
    'LOW',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP()
    1
);

INSERT INTO tasks (title, content, priority, created, finished, workspace_id) VALUES (
    'Second Task',
    'This is the content of the second task',
    'MEDIUM',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP()
    1
);

INSERT INTO tasks (title, content, priority, created, finished, workspace_id) VALUES (
    'Third Task',
    'This is the content of the third task',
    'HIGH',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP()
    1
);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
