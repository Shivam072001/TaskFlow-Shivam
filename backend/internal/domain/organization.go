package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrgRole string

const (
	OrgRoleOwner   OrgRole = "owner"
	OrgRoleAdmin   OrgRole = "admin"
	OrgRoleManager OrgRole = "manager"
	OrgRoleMember  OrgRole = "member"
)

func (r OrgRole) IsValid() bool {
	return r == OrgRoleOwner || r == OrgRoleAdmin || r == OrgRoleManager || r == OrgRoleMember
}

func (r OrgRole) CanManageOrg() bool      { return r == OrgRoleOwner || r == OrgRoleAdmin }
func (r OrgRole) CanInviteToOrg() bool    { return r == OrgRoleOwner || r == OrgRoleAdmin }
func (r OrgRole) CanManageTeams() bool    { return r == OrgRoleOwner || r == OrgRoleAdmin || r == OrgRoleManager }
func (r OrgRole) CanCreateWorkspace() bool { return r.IsValid() }

type Organization struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type OrgMember struct {
	ID        uuid.UUID `json:"id"`
	OrgID     uuid.UUID `json:"org_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      OrgRole   `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
	UserName  string    `json:"user_name,omitempty"`
	UserEmail string    `json:"user_email,omitempty"`
}
