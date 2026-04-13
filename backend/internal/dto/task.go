package dto

import (
	"strings"
	"time"
)

type CreateTaskRequest struct {
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Priority     string            `json:"priority"`
	AssigneeID   *string           `json:"assignee_id"`
	StartDate    *string           `json:"start_date"`
	DueDate      *string           `json:"due_date"`
	CustomFields map[string]string `json:"custom_fields"`
}

func (r *CreateTaskRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Title) == "" {
		errs["title"] = "is required"
	}
	if r.Priority == "" {
		r.Priority = "medium"
	}
	if r.Priority != "low" && r.Priority != "medium" && r.Priority != "high" {
		errs["priority"] = "must be low, medium, or high"
	}
	if r.StartDate != nil {
		if _, err := time.Parse("2006-01-02", *r.StartDate); err != nil {
			errs["start_date"] = "must be in YYYY-MM-DD format"
		}
	}
	if r.DueDate != nil {
		if _, err := time.Parse("2006-01-02", *r.DueDate); err != nil {
			errs["due_date"] = "must be in YYYY-MM-DD format"
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type UpdateTaskRequest struct {
	Title         *string           `json:"title"`
	Description   *string           `json:"description"`
	Status        *string           `json:"status"`
	Priority      *string           `json:"priority"`
	AssigneeID    *string           `json:"assignee_id"`
	StartDate     *string           `json:"start_date"`
	DueDate       *string           `json:"due_date"`
	BlockedReason *string           `json:"blocked_reason"`
	BlockedByTask *string           `json:"blocked_by_task"`
	CustomFields  map[string]string `json:"custom_fields"`
}

func (r *UpdateTaskRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.Title != nil && strings.TrimSpace(*r.Title) == "" {
		errs["title"] = "cannot be empty"
	}
	if r.Status != nil {
		s := *r.Status
		if s != "todo" && s != "in_progress" && s != "done" && s != "blocked" {
			errs["status"] = "must be todo, in_progress, done, or blocked"
		}
		if s == "blocked" {
			reason := ""
			byTask := ""
			if r.BlockedReason != nil {
				reason = strings.TrimSpace(*r.BlockedReason)
			}
			if r.BlockedByTask != nil {
				byTask = strings.TrimSpace(*r.BlockedByTask)
			}
			if reason == "" && byTask == "" {
				errs["blocked_reason"] = "either blocked_reason or blocked_by_task is required"
			}
		}
	}
	if r.Priority != nil {
		p := *r.Priority
		if p != "low" && p != "medium" && p != "high" {
			errs["priority"] = "must be low, medium, or high"
		}
	}
	if r.StartDate != nil && *r.StartDate != "" {
		if _, err := time.Parse("2006-01-02", *r.StartDate); err != nil {
			errs["start_date"] = "must be in YYYY-MM-DD format"
		}
	}
	if r.DueDate != nil && *r.DueDate != "" {
		if _, err := time.Parse("2006-01-02", *r.DueDate); err != nil {
			errs["due_date"] = "must be in YYYY-MM-DD format"
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
