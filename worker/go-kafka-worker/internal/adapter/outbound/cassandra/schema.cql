-- Local dev reference schema for JobStore. Not run automatically.
CREATE TABLE IF NOT EXISTS jobs (
    id           text PRIMARY KEY,
    task_id      text,
    type         text,
    payload      text,
    received_at  timestamp,
    processed_at timestamp
);
