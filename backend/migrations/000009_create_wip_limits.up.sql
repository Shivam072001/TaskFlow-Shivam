CREATE TABLE project_wip_limits (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    status     TEXT NOT NULL CHECK (status IN ('todo', 'in_progress', 'done')),
    max_tasks  INT NOT NULL CHECK (max_tasks > 0),
    UNIQUE (project_id, status)
);
