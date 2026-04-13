package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type ProjectMemberService struct {
	projectMemberRepo *repository.ProjectMemberRepository
	workspaceRepo     *repository.WorkspaceRepository
}

func NewProjectMemberService(pmr *repository.ProjectMemberRepository, wr *repository.WorkspaceRepository) *ProjectMemberService {
	return &ProjectMemberService{projectMemberRepo: pmr, workspaceRepo: wr}
}

func (s *ProjectMemberService) AddMember(ctx context.Context, projectID, workspaceID, userID uuid.UUID, role domain.WorkspaceRole) (*domain.ProjectMember, error) {
	// Validate user is a workspace member
	_, err := s.workspaceRepo.GetMember(ctx, workspaceID, userID)
	if err != nil {
		return nil, fmt.Errorf("user is not a workspace member: %w", domain.ErrForbidden)
	}

	now := time.Now().UTC()
	m := &domain.ProjectMember{
		ID:        uuid.New(),
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		JoinedAt:  now,
	}
	if err := s.projectMemberRepo.Add(ctx, m); err != nil {
		return nil, fmt.Errorf("adding project member: %w", err)
	}
	return m, nil
}

func (s *ProjectMemberService) RemoveMember(ctx context.Context, projectID, userID uuid.UUID) error {
	return s.projectMemberRepo.Remove(ctx, projectID, userID)
}

func (s *ProjectMemberService) ListMembers(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectMember, error) {
	return s.projectMemberRepo.List(ctx, projectID)
}

func (s *ProjectMemberService) GetMember(ctx context.Context, projectID, userID uuid.UUID) (*domain.ProjectMember, error) {
	return s.projectMemberRepo.GetMember(ctx, projectID, userID)
}

func (s *ProjectMemberService) UpdateRole(ctx context.Context, projectID, userID uuid.UUID, role domain.WorkspaceRole) error {
	return s.projectMemberRepo.UpdateRole(ctx, projectID, userID, role)
}
