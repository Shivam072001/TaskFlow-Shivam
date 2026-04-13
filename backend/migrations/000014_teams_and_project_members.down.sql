DROP TABLE IF EXISTS project_teams;
DROP TABLE IF EXISTS workspace_teams;
DROP TABLE IF EXISTS project_members;
DROP INDEX IF EXISTS idx_org_inv_unique_pending;
DROP TABLE IF EXISTS org_invitations;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;

ALTER TABLE organization_members DROP CONSTRAINT IF EXISTS organization_members_role_check;
ALTER TABLE organization_members ADD CONSTRAINT organization_members_role_check
  CHECK (role IN ('owner', 'admin', 'member'));
