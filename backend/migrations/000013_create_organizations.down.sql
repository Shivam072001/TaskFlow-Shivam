ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_prefix_org_unique;
ALTER TABLE projects ADD CONSTRAINT projects_prefix_workspace_unique UNIQUE (workspace_id, prefix);
ALTER TABLE projects DROP COLUMN IF EXISTS org_id;
ALTER TABLE workspaces DROP COLUMN IF EXISTS org_id;
DROP TABLE IF EXISTS organization_members;
DROP TABLE IF EXISTS organizations;
