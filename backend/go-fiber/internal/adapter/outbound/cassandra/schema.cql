-- Local dev reference schema for TaskRepository. Not run automatically.
CREATE TABLE IF NOT EXISTS tasks (
    id         text PRIMARY KEY,
    title      text,
    done       boolean,
    created_at timestamp,
    updated_at timestamp
);
