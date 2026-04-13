package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceRole string

const (
	RoleOwner     WorkspaceRole = "owner"
	RoleAdmin     WorkspaceRole = "admin"
	RoleManager   WorkspaceRole = "manager"
	RoleLead      WorkspaceRole = "lead"
	RoleMember    WorkspaceRole = "member"
	RoleGuest     WorkspaceRole = "guest"
	RoleViewer    WorkspaceRole = "viewer"
)

// RolePower returns a numeric level for hierarchy comparisons.
// Higher number = more power.
var RolePower = map[WorkspaceRole]int{
	RoleOwner:   70,
	RoleAdmin:   60,
	RoleManager: 50,
	RoleLead:    40,
	RoleMember:  30,
	RoleGuest:   20,
	RoleViewer:  10,
}

func (r WorkspaceRole) IsValid() bool {
	_, ok := RolePower[r]
	return ok
}

func (r WorkspaceRole) IsAssignable() bool {
	return r != RoleOwner && r.IsValid()
}

func (r WorkspaceRole) Power() int {
	return RolePower[r]
}

func (r WorkspaceRole) IsSuperiorTo(other WorkspaceRole) bool {
	return RolePower[r] > RolePower[other]
}

// Permission checks
func (r WorkspaceRole) CanManageWorkspace() bool  { return r == RoleOwner || r == RoleAdmin }
func (r WorkspaceRole) CanInviteMembers() bool     { return RolePower[r] >= RolePower[RoleManager] }
func (r WorkspaceRole) CanManageMembers() bool     { return RolePower[r] >= RolePower[RoleAdmin] }
func (r WorkspaceRole) CanCreateProject() bool     { return RolePower[r] >= RolePower[RoleMember] }
func (r WorkspaceRole) CanEditAnyProject() bool    { return RolePower[r] >= RolePower[RoleManager] }
func (r WorkspaceRole) CanDeleteProject() bool     { return RolePower[r] >= RolePower[RoleManager] }
func (r WorkspaceRole) CanCreateTask() bool        { return RolePower[r] >= RolePower[RoleMember] }
func (r WorkspaceRole) CanEditAnyTask() bool       { return RolePower[r] >= RolePower[RoleLead] }
func (r WorkspaceRole) CanDeleteAnyTask() bool     { return RolePower[r] >= RolePower[RoleManager] }
func (r WorkspaceRole) CanComment() bool           { return RolePower[r] >= RolePower[RoleGuest] }
func (r WorkspaceRole) CanManageCustomFields() bool { return RolePower[r] >= RolePower[RoleAdmin] }
func (r WorkspaceRole) CanSetWIPLimits() bool      { return RolePower[r] >= RolePower[RoleAdmin] }
func (r WorkspaceRole) CanView() bool              { return r.IsValid() }

type Workspace struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WorkspaceMember struct {
	ID          uuid.UUID     `json:"id"`
	WorkspaceID uuid.UUID     `json:"workspace_id"`
	UserID      uuid.UUID     `json:"user_id"`
	Role        WorkspaceRole `json:"role"`
	JoinedAt    time.Time     `json:"joined_at"`
	UserName    string        `json:"user_name,omitempty"`
	UserEmail   string        `json:"user_email,omitempty"`
}

type InvitationStatus string

const (
	InvitePending  InvitationStatus = "pending"
	InviteAccepted InvitationStatus = "accepted"
	InviteDeclined InvitationStatus = "declined"
)

type WorkspaceInvitation struct {
	ID            uuid.UUID        `json:"id"`
	WorkspaceID   uuid.UUID        `json:"workspace_id"`
	InviterID     uuid.UUID        `json:"inviter_id"`
	InviteeEmail  string           `json:"invitee_email"`
	InviteeID     *uuid.UUID       `json:"invitee_id"`
	Role          WorkspaceRole    `json:"role"`
	Status        InvitationStatus `json:"status"`
	CreatedAt     time.Time        `json:"created_at"`
	RespondedAt   *time.Time       `json:"responded_at"`
	WorkspaceName string           `json:"workspace_name,omitempty"`
	InviterName   string           `json:"inviter_name,omitempty"`
}
