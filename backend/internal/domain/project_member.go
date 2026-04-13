package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMember struct {
	ID        uuid.UUID     `json:"id"`
	ProjectID uuid.UUID     `json:"project_id"`
	UserID    uuid.UUID     `json:"user_id"`
	Role      WorkspaceRole `json:"role"`
	JoinedAt  time.Time     `json:"joined_at"`
	UserName  string        `json:"user_name,omitempty"`
	UserEmail string        `json:"user_email,omitempty"`
}
