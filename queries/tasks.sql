CREATE TABLE tasks (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    priority TEXT NOT NULL,
    created DATETIME NOT NULL,
    finished DATETIME
);

CREATE INDEX idx_tasks_created ON tasks(created);

INSERT INTO tasks (title, content, priority, created, finished) VALUES (
    'First Task',
    'This is the content of the first task',
    'LOW',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP()
);

INSERT INTO tasks (title, content, priority, created, finished) VALUES (
    'Second Task',
    'This is the content of the second task',
    'MEDIUM',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP()
);

INSERT INTO tasks (title, content, priority, created, finished) VALUES (
    'Third Task',
    'This is the content of the third task',
    'HIGH',
    UTC_TIMESTAMP(),
    UTC_TIMESTAMP()
);
