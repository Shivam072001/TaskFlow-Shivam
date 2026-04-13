package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
}

func NewCommentService(cr *repository.CommentRepository) *CommentService {
	return &CommentService{commentRepo: cr}
}

func (s *CommentService) Create(ctx context.Context, entityType domain.CommentEntityType, entityID, userID uuid.UUID, parentID *uuid.UUID, content string) (*domain.Comment, error) {
	now := time.Now().UTC()
	c := &domain.Comment{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		UserID:     userID,
		ParentID:   parentID,
		Content:    content,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.commentRepo.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("creating comment: %w", err)
	}
	return c, nil
}

func (s *CommentService) List(ctx context.Context, entityType domain.CommentEntityType, entityID uuid.UUID) ([]domain.Comment, error) {
	return s.commentRepo.ListByEntity(ctx, entityType, entityID)
}

func (s *CommentService) Update(ctx context.Context, commentID, userID uuid.UUID, content string) (*domain.Comment, error) {
	c, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if c.UserID != userID {
		return nil, domain.ErrForbidden
	}

	c.Content = content
	c.UpdatedAt = time.Now().UTC()

	if err := s.commentRepo.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("updating comment: %w", err)
	}
	return c, nil
}

func (s *CommentService) Delete(ctx context.Context, commentID, userID uuid.UUID, callerRole domain.WorkspaceRole) error {
	c, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	isOwner := c.UserID == userID
	canModerate := callerRole.Power() >= domain.RolePower[domain.RoleManager]

	if !isOwner && !canModerate {
		return domain.ErrForbidden
	}

	return s.commentRepo.Delete(ctx, commentID)
}
