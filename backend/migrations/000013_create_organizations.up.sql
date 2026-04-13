-- Organizations
CREATE TABLE organizations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    created_by  UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_organizations_slug ON organizations (slug);

-- Org membership
CREATE TABLE organization_members (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id    UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role      TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member')),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (org_id, user_id)
);

CREATE INDEX idx_org_members_org ON organization_members (org_id);
CREATE INDEX idx_org_members_user ON organization_members (user_id);

-- Scope workspaces under orgs
ALTER TABLE workspaces ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

-- Backfill: create a personal org per workspace creator, assign workspaces
DO $$
DECLARE
    r RECORD;
    org_uuid UUID;
BEGIN
    FOR r IN SELECT DISTINCT created_by FROM workspaces LOOP
        org_uuid := gen_random_uuid();
        INSERT INTO organizations (id, name, slug, created_by, created_at)
        SELECT org_uuid,
               u.name || '''s Organization',
               LOWER(REPLACE(u.name, ' ', '-')) || '-' || LEFT(org_uuid::text, 8),
               r.created_by,
               now()
        FROM users u WHERE u.id = r.created_by;

        INSERT INTO organization_members (id, org_id, user_id, role, joined_at)
        VALUES (gen_random_uuid(), org_uuid, r.created_by, 'owner', now());

        UPDATE workspaces SET org_id = org_uuid WHERE created_by = r.created_by;
    END LOOP;
END $$;

ALTER TABLE workspaces ALTER COLUMN org_id SET NOT NULL;

-- Scope projects under orgs (derived from workspace)
ALTER TABLE projects ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

UPDATE projects p SET org_id = w.org_id FROM workspaces w WHERE w.id = p.workspace_id;

ALTER TABLE projects ALTER COLUMN org_id SET NOT NULL;

-- Move prefix uniqueness from workspace-scoped to org-scoped
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_prefix_workspace_unique;
ALTER TABLE projects ADD CONSTRAINT projects_prefix_org_unique UNIQUE (org_id, prefix);

-- Add org members for all workspace members who aren't already org members
INSERT INTO organization_members (id, org_id, user_id, role, joined_at)
SELECT gen_random_uuid(), w.org_id, wm.user_id, 'member', wm.joined_at
FROM workspace_members wm
JOIN workspaces w ON w.id = wm.workspace_id
WHERE NOT EXISTS (
    SELECT 1 FROM organization_members om
    WHERE om.org_id = w.org_id AND om.user_id = wm.user_id
);
