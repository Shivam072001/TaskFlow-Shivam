package domain

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Name        string    `json:"name"`
	Prefix      string    `json:"prefix"`
	Description string    `json:"description"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type WIPLimit struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"project_id"`
	Status    TaskStatus `json:"status"`
	MaxTasks  int        `json:"max_tasks"`
}

type CustomFieldType string

const (
	FieldTypeText   CustomFieldType = "text"
	FieldTypeNumber CustomFieldType = "number"
	FieldTypeSelect CustomFieldType = "select"
)

type CustomFieldDefinition struct {
	ID        uuid.UUID       `json:"id"`
	ProjectID uuid.UUID       `json:"project_id"`
	Name      string          `json:"name"`
	FieldType CustomFieldType `json:"field_type"`
	Options   []string        `json:"options,omitempty"`
	Required  bool            `json:"required"`
	CreatedBy uuid.UUID       `json:"created_by"`
	CreatedAt time.Time       `json:"created_at"`
}

type CustomFieldValue struct {
	ID      uuid.UUID `json:"id"`
	TaskID  uuid.UUID `json:"task_id"`
	FieldID uuid.UUID `json:"field_id"`
	Value   string    `json:"value"`
}
