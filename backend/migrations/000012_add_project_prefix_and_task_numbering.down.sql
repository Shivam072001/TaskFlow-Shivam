ALTER TABLE tasks DROP COLUMN IF EXISTS blocked_by_task;
ALTER TABLE tasks DROP COLUMN IF EXISTS blocked_reason;

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_status_check;
ALTER TABLE tasks ADD CONSTRAINT tasks_status_check CHECK (status IN ('todo', 'in_progress', 'done'));

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_number_project_unique;
ALTER TABLE tasks DROP COLUMN IF EXISTS task_key;
ALTER TABLE tasks DROP COLUMN IF EXISTS task_number;

ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_prefix_workspace_unique;
ALTER TABLE projects DROP COLUMN IF EXISTS prefix;
