package dto

import "strings"

type CreateWorkspaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *CreateWorkspaceRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Name) == "" {
		errs["name"] = "is required"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdateWorkspaceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (r *UpdateWorkspaceRequest) Validate() map[string]string {
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return map[string]string{"name": "cannot be empty"}
	}
	return nil
}

var assignableRoles = map[string]bool{
	"admin": true, "manager": true, "lead": true,
	"member": true, "guest": true, "viewer": true,
}

type InviteMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (r *InviteMemberRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = "is required"
	}
	if r.Role == "" {
		r.Role = "member"
	}
	if !assignableRoles[r.Role] {
		errs["role"] = "must be one of: admin, manager, lead, member, guest, viewer"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type InvitationResponseRequest struct {
	Action string `json:"action"`
}

func (r *InvitationResponseRequest) Validate() map[string]string {
	if r.Action != "accept" && r.Action != "decline" {
		return map[string]string{"action": "must be accept or decline"}
	}
	return nil
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role"`
}

func (r *UpdateMemberRoleRequest) Validate() map[string]string {
	if !assignableRoles[r.Role] {
		return map[string]string{"role": "must be one of: admin, manager, lead, member, guest, viewer"}
	}
	return nil
}

type WorkspaceStatsResponse struct {
	ProjectCount    int               `json:"project_count"`
	TasksByStatus   map[string]int    `json:"tasks_by_status"`
	TasksByAssignee []AssigneeStats   `json:"tasks_by_assignee"`
	OverdueCount    int               `json:"overdue_count"`
}

type AssigneeStats struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Count    int    `json:"count"`
}
