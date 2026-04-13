package domain

import (
	"time"

	"github.com/google/uuid"
)

type CommentEntityType string

const (
	EntityProject CommentEntityType = "project"
	EntityTask    CommentEntityType = "task"
)

type Comment struct {
	ID         uuid.UUID         `json:"id"`
	EntityType CommentEntityType `json:"entity_type"`
	EntityID   uuid.UUID         `json:"entity_id"`
	UserID     uuid.UUID         `json:"user_id"`
	ParentID   *uuid.UUID        `json:"parent_id"`
	Content    string            `json:"content"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	UserName   string            `json:"user_name,omitempty"`
	UserEmail  string            `json:"user_email,omitempty"`
	Replies    []Comment         `json:"replies,omitempty"`
}
