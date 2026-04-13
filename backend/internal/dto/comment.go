package dto

import "strings"

type CreateCommentRequest struct {
	EntityType string  `json:"entity_type"`
	EntityID   string  `json:"entity_id"`
	ParentID   *string `json:"parent_id"`
	Content    string  `json:"content"`
}

func (r *CreateCommentRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.EntityType != "project" && r.EntityType != "task" {
		errs["entity_type"] = "must be project or task"
	}
	if strings.TrimSpace(r.EntityID) == "" {
		errs["entity_id"] = "is required"
	}
	if strings.TrimSpace(r.Content) == "" {
		errs["content"] = "is required"
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

func (r *UpdateCommentRequest) Validate() map[string]string {
	if strings.TrimSpace(r.Content) == "" {
		return map[string]string{"content": "is required"}
	}
	return nil
}
