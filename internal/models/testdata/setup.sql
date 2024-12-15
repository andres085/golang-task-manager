CREATE TABLE users (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	firstName VARCHAR(255) NOT NULL,
	lastName VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (firstName, lastName, email, hashed_password, created) VALUES (
    'Test',
    'McTester',
    'test@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 09:18:24'
);

INSERT INTO users (firstName, lastName, email, hashed_password, created) VALUES (
    'Member',
    'Memberino',
    'member@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 09:18:24'
);

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
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_id ON users_workspaces (user_id);
CREATE INDEX idx_workspace_id ON users_workspaces (workspace_id);

INSERT INTO workspaces (title, description, created) VALUES (
    'First Workspace',
    'This is the first workspace description',
    UTC_TIMESTAMP()
);

CREATE TABLE tasks (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    priority TEXT NOT NULL,
    created DATETIME NOT NULL,
    finished DATETIME DEFAULT NULL,
    workspace_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'To Do',
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO users_workspaces(user_id, workspace_id, role, created) VALUES (1, 1, "ADMIN", UTC_TIMESTAMP());
INSERT INTO users_workspaces(user_id, workspace_id, role, created) VALUES (2, 1, "MEMBER", UTC_TIMESTAMP());

CREATE INDEX idx_tasks_created ON tasks(created);

INSERT INTO tasks (title, content, priority, created, finished, workspace_id, user_id) VALUES (
    'First Task',
    'This is the content of the first task',
    'LOW',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP(),
    1,
    1
);

INSERT INTO tasks (title, content, priority, created, finished, workspace_id, user_id) VALUES (
    'Second Task',
    'This is the content of the second task',
    'MEDIUM',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP(),
    1,
    1
);

INSERT INTO tasks (title, content, priority, created, finished, workspace_id, user_id) VALUES (
    'Third Task',
    'This is the content of the third task',
    'HIGH',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP(),
    1,
    1
);
