ALTER TABLE projects ADD COLUMN prefix TEXT NOT NULL DEFAULT '';
ALTER TABLE projects ADD CONSTRAINT projects_prefix_workspace_unique UNIQUE (workspace_id, prefix);

ALTER TABLE tasks ADD COLUMN task_number INT;
ALTER TABLE tasks ADD COLUMN task_key TEXT;
ALTER TABLE tasks ADD CONSTRAINT tasks_number_project_unique UNIQUE (project_id, task_number);

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_status_check;
ALTER TABLE tasks ADD CONSTRAINT tasks_status_check CHECK (status IN ('todo', 'in_progress', 'done', 'blocked'));

ALTER TABLE tasks ADD COLUMN blocked_reason TEXT NOT NULL DEFAULT '';
ALTER TABLE tasks ADD COLUMN blocked_by_task TEXT NOT NULL DEFAULT '';
