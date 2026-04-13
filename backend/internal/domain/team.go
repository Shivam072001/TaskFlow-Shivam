package domain

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `json:"id"`
	OrgID     uuid.UUID `json:"org_id"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type TeamMember struct {
	ID        uuid.UUID `json:"id"`
	TeamID    uuid.UUID `json:"team_id"`
	UserID    uuid.UUID `json:"user_id"`
	AddedAt   time.Time `json:"added_at"`
	UserName  string    `json:"user_name,omitempty"`
	UserEmail string    `json:"user_email,omitempty"`
}

type WorkspaceTeam struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	TeamID      uuid.UUID `json:"team_id"`
	DefaultRole string    `json:"default_role"`
	AddedAt     time.Time `json:"added_at"`
}

type ProjectTeam struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	TeamID      uuid.UUID `json:"team_id"`
	DefaultRole string    `json:"default_role"`
	AddedAt     time.Time `json:"added_at"`
}
