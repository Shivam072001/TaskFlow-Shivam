-- Teams (org-scoped)
CREATE TABLE teams (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    created_by  UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (org_id, name)
);

-- Team members (org members assigned to a team)
CREATE TABLE team_members (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id   UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    added_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (team_id, user_id)
);

-- Expand org_members role CHECK to include 'manager'
ALTER TABLE organization_members DROP CONSTRAINT IF EXISTS organization_members_role_check;
ALTER TABLE organization_members ADD CONSTRAINT organization_members_role_check
  CHECK (role IN ('owner', 'admin', 'manager', 'member'));

-- Org-level invitations
CREATE TABLE org_invitations (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id         UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    inviter_id     UUID NOT NULL REFERENCES users(id),
    invitee_email  TEXT NOT NULL,
    invitee_id     UUID REFERENCES users(id),
    role           TEXT NOT NULL CHECK (role IN ('admin', 'manager', 'member')),
    status         TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    responded_at   TIMESTAMPTZ
);
CREATE UNIQUE INDEX idx_org_inv_unique_pending ON org_invitations (org_id, invitee_email) WHERE status = 'pending';

-- Project members
CREATE TABLE project_members (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role        TEXT NOT NULL CHECK (role IN ('owner','admin','manager','lead','member','guest','viewer')),
    joined_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, user_id)
);

-- Track which teams are assigned to workspaces (for auto-cascade)
CREATE TABLE workspace_teams (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id  UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    team_id       UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    default_role  TEXT NOT NULL DEFAULT 'member',
    added_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (workspace_id, team_id)
);

-- Track which teams are assigned to projects (for auto-cascade)
CREATE TABLE project_teams (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id    UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    team_id       UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    default_role  TEXT NOT NULL DEFAULT 'member',
    added_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (project_id, team_id)
);
