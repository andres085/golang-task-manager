CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL,
    INDEX sessions_expiry_idx (expiry)
);

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    firstName VARCHAR(255) NOT NULL,
    lastName VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    UNIQUE KEY users_uc_email (email)
);

CREATE TABLE workspaces (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created DATETIME NOT NULL
);

CREATE TABLE users_workspaces (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    workspace_id INT NOT NULL,
    role TEXT NOT NULL,
    created DATETIME NOT NULL,
    UNIQUE KEY users_workspaces_uc (user_id, workspace_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_workspace_id (workspace_id)
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
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_tasks_created (created)
);
