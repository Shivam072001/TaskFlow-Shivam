CREATE TABLE workspace_invitations (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id   UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    inviter_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    invitee_email  TEXT NOT NULL,
    invitee_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    role           TEXT NOT NULL CHECK (role IN ('admin', 'manager', 'lead', 'member', 'guest', 'viewer')),
    status         TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    responded_at   TIMESTAMPTZ
);

CREATE INDEX idx_invitations_workspace ON workspace_invitations (workspace_id);
CREATE INDEX idx_invitations_invitee ON workspace_invitations (invitee_email);
CREATE INDEX idx_invitations_invitee_id ON workspace_invitations (invitee_id);
CREATE UNIQUE INDEX idx_invitations_unique_pending ON workspace_invitations (workspace_id, invitee_email) WHERE status = 'pending';
