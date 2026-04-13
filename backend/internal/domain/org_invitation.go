package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrgInvitation struct {
	ID           uuid.UUID  `json:"id"`
	OrgID        uuid.UUID  `json:"org_id"`
	InviterID    uuid.UUID  `json:"inviter_id"`
	InviteeEmail string     `json:"invitee_email"`
	InviteeID    *uuid.UUID `json:"invitee_id"`
	Role         OrgRole    `json:"role"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	RespondedAt  *time.Time `json:"responded_at"`
	OrgName      string     `json:"org_name,omitempty"`
	InviterName  string     `json:"inviter_name,omitempty"`
}
