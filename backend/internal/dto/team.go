package dto

import "strings"

type CreateTeamRequest struct {
	Name string `json:"name"`
}

func (r *CreateTeamRequest) Validate() map[string]string {
	if strings.TrimSpace(r.Name) == "" {
		return map[string]string{"name": "is required"}
	}
	return nil
}

type UpdateTeamRequest struct {
	Name string `json:"name"`
}

func (r *UpdateTeamRequest) Validate() map[string]string {
	if strings.TrimSpace(r.Name) == "" {
		return map[string]string{"name": "is required"}
	}
	return nil
}

type AddTeamMemberRequest struct {
	UserID string `json:"user_id"`
}

func (r *AddTeamMemberRequest) Validate() map[string]string {
	if strings.TrimSpace(r.UserID) == "" {
		return map[string]string{"user_id": "is required"}
	}
	return nil
}

type AddTeamToEntityRequest struct {
	TeamID      string `json:"team_id"`
	DefaultRole string `json:"default_role"`
}

func (r *AddTeamToEntityRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.TeamID) == "" {
		errs["team_id"] = "is required"
	}
	if strings.TrimSpace(r.DefaultRole) == "" {
		errs["default_role"] = "is required"
	} else if !assignableRoles[r.DefaultRole] {
		errs["default_role"] = "must be one of: admin, manager, lead, member, guest, viewer"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type DirectAddMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func (r *DirectAddMemberRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.UserID) == "" {
		errs["user_id"] = "is required"
	}
	if strings.TrimSpace(r.Role) == "" {
		errs["role"] = "is required"
	} else if !assignableRoles[r.Role] {
		errs["role"] = "must be one of: admin, manager, lead, member, guest, viewer"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
