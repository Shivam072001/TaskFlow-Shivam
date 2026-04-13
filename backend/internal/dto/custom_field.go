package dto

import "strings"

type CreateCustomFieldRequest struct {
	Name      string   `json:"name"`
	FieldType string   `json:"field_type"`
	Options   []string `json:"options"`
	Required  bool     `json:"required"`
}

func (r *CreateCustomFieldRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Name) == "" {
		errs["name"] = "is required"
	}
	if r.FieldType != "text" && r.FieldType != "number" && r.FieldType != "select" {
		errs["field_type"] = "must be text, number, or select"
	}
	if r.FieldType == "select" && len(r.Options) == 0 {
		errs["options"] = "required for select fields"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdateCustomFieldRequest struct {
	Name     *string  `json:"name"`
	Options  []string `json:"options"`
	Required *bool    `json:"required"`
}

func (r *UpdateCustomFieldRequest) Validate() map[string]string {
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return map[string]string{"name": "cannot be empty"}
	}
	return nil
}

type SetFieldValueRequest struct {
	Value string `json:"value"`
}

func (r *SetFieldValueRequest) Validate() map[string]string {
	return nil
}

type SetWIPLimitRequest struct {
	Status   string `json:"status"`
	MaxTasks int    `json:"max_tasks"`
}

func (r *SetWIPLimitRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Status != "todo" && r.Status != "in_progress" && r.Status != "done" {
		errs["status"] = "must be todo, in_progress, or done"
	}
	if r.MaxTasks < 1 {
		errs["max_tasks"] = "must be at least 1"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
