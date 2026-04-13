package dto

import (
	"regexp"
	"strings"
)

var prefixPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]{1,5}$`)

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
}

func (r *CreateProjectRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Name) == "" {
		errs["name"] = "is required"
	}
	r.Prefix = strings.TrimSpace(strings.ToUpper(r.Prefix))
	if r.Prefix == "" {
		errs["prefix"] = "is required"
	} else if !prefixPattern.MatchString(r.Prefix) {
		errs["prefix"] = "must be 2-6 uppercase alphanumeric characters starting with a letter"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (r *UpdateProjectRequest) Validate() map[string]string {
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return map[string]string{"name": "cannot be empty"}
	}
	return nil
}
