-- Local dev reference schema for JobStore. Not run automatically — a
-- real project wires this into its migration tool of choice.
CREATE TABLE IF NOT EXISTS jobs (
    id           TEXT PRIMARY KEY,
    task_id      TEXT NOT NULL,
    type         TEXT NOT NULL,
    payload      TEXT NOT NULL,
    received_at  TIMESTAMPTZ NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL
);
