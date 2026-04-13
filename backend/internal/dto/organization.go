package dto

import (
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{1,28}[a-z0-9]$`)

type CreateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (r *CreateOrganizationRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Name) == "" {
		errs["name"] = "is required"
	}
	r.Slug = strings.TrimSpace(strings.ToLower(r.Slug))
	if r.Slug == "" {
		errs["slug"] = "is required"
	} else if !slugPattern.MatchString(r.Slug) {
		errs["slug"] = "must be 3-30 lowercase alphanumeric characters or hyphens, starting with a letter"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdateOrganizationRequest struct {
	Name *string `json:"name"`
}

func (r *UpdateOrganizationRequest) Validate() map[string]string {
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return map[string]string{"name": "cannot be empty"}
	}
	return nil
}
