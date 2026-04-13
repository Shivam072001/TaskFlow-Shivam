package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
	StatusBlocked    TaskStatus = "blocked"
)

type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID            uuid.UUID         `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Status        TaskStatus        `json:"status"`
	Priority      TaskPriority      `json:"priority"`
	ProjectID     uuid.UUID         `json:"project_id"`
	AssigneeID    *uuid.UUID        `json:"assignee_id"`
	StartDate     *time.Time        `json:"start_date"`
	DueDate       *time.Time        `json:"due_date"`
	CreatedBy     uuid.UUID         `json:"created_by"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	TaskNumber    int               `json:"task_number"`
	TaskKey       string            `json:"task_key"`
	BlockedReason string            `json:"blocked_reason"`
	BlockedByTask string            `json:"blocked_by_task"`
	CustomFields  map[string]string `json:"custom_fields,omitempty"`
}
