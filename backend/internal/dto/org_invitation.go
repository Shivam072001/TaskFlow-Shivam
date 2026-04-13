package dto

import "strings"

var validOrgInviteRoles = map[string]bool{
	"admin": true, "manager": true, "member": true,
}

type SendOrgInviteRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (r *SendOrgInviteRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = "is required"
	}
	if r.Role == "" {
		r.Role = "member"
	}
	if !validOrgInviteRoles[r.Role] {
		errs["role"] = "must be one of: admin, manager, member"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type RespondOrgInviteRequest struct {
	Action string `json:"action"`
}

func (r *RespondOrgInviteRequest) Validate() map[string]string {
	if r.Action != "accept" && r.Action != "decline" {
		return map[string]string{"action": "must be accept or decline"}
	}
	return nil
}
