-- Local dev reference schema for CassandraTaskRepository. Not run
-- automatically.
CREATE TABLE IF NOT EXISTS tasks (
    id         text PRIMARY KEY,
    title      text,
    done       boolean,
    created_at timestamp,
    updated_at timestamp
);
