package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type ProjectService struct {
	projectRepo   *repository.ProjectRepository
	workspaceRepo *repository.WorkspaceRepository
}

func NewProjectService(pr *repository.ProjectRepository, wr *repository.WorkspaceRepository) *ProjectService {
	return &ProjectService{projectRepo: pr, workspaceRepo: wr}
}

func (s *ProjectService) Create(ctx context.Context, workspaceID, ownerID uuid.UUID, name, prefix, description string) (*domain.Project, error) {
	ws, err := s.workspaceRepo.GetByID(ctx, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("fetching workspace for org_id: %w", err)
	}

	p := &domain.Project{
		ID:          uuid.New(),
		OrgID:       ws.OrgID,
		Name:        name,
		Prefix:      prefix,
		Description: description,
		WorkspaceID: workspaceID,
		OwnerID:     ownerID,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.projectRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("creating project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

func (s *ProjectService) ListByWorkspace(ctx context.Context, f repository.ProjectFilter) ([]domain.Project, int, error) {
	return s.projectRepo.ListByWorkspace(ctx, f)
}

func (s *ProjectService) Update(ctx context.Context, id, callerID uuid.UUID, name, description *string) (*domain.Project, error) {
	p, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.OwnerID != callerID {
		return nil, domain.ErrForbidden
	}

	if name != nil {
		p.Name = *name
	}
	if description != nil {
		p.Description = *description
	}

	if err := s.projectRepo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("updating project: %w", err)
	}
	return p, nil
}

func (s *ProjectService) Delete(ctx context.Context, id, callerID uuid.UUID) error {
	p, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.OwnerID != callerID {
		return domain.ErrForbidden
	}

	return s.projectRepo.Delete(ctx, id)
}

// VerifyWorkspaceMembership checks that the caller is a member of the project's workspace
func (s *ProjectService) VerifyWorkspaceMembership(ctx context.Context, projectID, userID uuid.UUID) (*domain.Project, error) {
	p, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	_, err = s.workspaceRepo.GetMember(ctx, p.WorkspaceID, userID)
	if err != nil {
		return nil, domain.ErrForbidden
	}

	return p, nil
}
