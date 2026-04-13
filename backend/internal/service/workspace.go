package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/dto"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type WorkspaceService struct {
	workspaceRepo     *repository.WorkspaceRepository
	userRepo          *repository.UserRepository
	taskRepo          *repository.TaskRepository
	projectMemberRepo *repository.ProjectMemberRepository
}

func NewWorkspaceService(wr *repository.WorkspaceRepository, ur *repository.UserRepository, tr *repository.TaskRepository, pmr *repository.ProjectMemberRepository) *WorkspaceService {
	return &WorkspaceService{workspaceRepo: wr, userRepo: ur, taskRepo: tr, projectMemberRepo: pmr}
}

func (s *WorkspaceService) Create(ctx context.Context, orgID, userID uuid.UUID, name, description string) (*domain.Workspace, error) {
	now := time.Now().UTC()
	w := &domain.Workspace{
		ID:          uuid.New(),
		OrgID:       orgID,
		Name:        name,
		Description: description,
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.workspaceRepo.Create(ctx, w); err != nil {
		return nil, fmt.Errorf("creating workspace: %w", err)
	}

	member := &domain.WorkspaceMember{
		ID:          uuid.New(),
		WorkspaceID: w.ID,
		UserID:      userID,
		Role:        domain.RoleOwner,
		JoinedAt:    now,
	}
	if err := s.workspaceRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("adding owner as member: %w", err)
	}

	return w, nil
}

func (s *WorkspaceService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Workspace, error) {
	return s.workspaceRepo.GetByID(ctx, id)
}

func (s *WorkspaceService) ListByUser(ctx context.Context, userID, orgID uuid.UUID, page, limit int) ([]domain.Workspace, int, error) {
	return s.workspaceRepo.ListByUser(ctx, userID, orgID, page, limit)
}

func (s *WorkspaceService) Update(ctx context.Context, id uuid.UUID, role domain.WorkspaceRole, name, description *string) (*domain.Workspace, error) {
	if !role.CanManageWorkspace() {
		return nil, domain.ErrForbidden
	}

	w, err := s.workspaceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		w.Name = *name
	}
	if description != nil {
		w.Description = *description
	}
	w.UpdatedAt = time.Now().UTC()

	if err := s.workspaceRepo.Update(ctx, w); err != nil {
		return nil, fmt.Errorf("updating workspace: %w", err)
	}
	return w, nil
}

func (s *WorkspaceService) Delete(ctx context.Context, id uuid.UUID, role domain.WorkspaceRole) error {
	if role != domain.RoleOwner {
		return domain.ErrForbidden
	}
	return s.workspaceRepo.Delete(ctx, id)
}

func (s *WorkspaceService) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]domain.WorkspaceMember, error) {
	return s.workspaceRepo.ListMembers(ctx, workspaceID)
}

func (s *WorkspaceService) UpdateMemberRole(ctx context.Context, workspaceID, targetUserID uuid.UUID, callerRole domain.WorkspaceRole, newRole domain.WorkspaceRole) error {
	if !callerRole.CanManageMembers() {
		return domain.ErrForbidden
	}
	if !callerRole.IsSuperiorTo(newRole) {
		return domain.ErrForbidden
	}

	target, err := s.workspaceRepo.GetMember(ctx, workspaceID, targetUserID)
	if err != nil {
		return err
	}
	if !callerRole.IsSuperiorTo(target.Role) {
		return domain.ErrForbidden
	}

	return s.workspaceRepo.UpdateMemberRole(ctx, workspaceID, targetUserID, newRole)
}

func (s *WorkspaceService) RemoveMember(ctx context.Context, workspaceID, targetUserID, callerUserID uuid.UUID, callerRole domain.WorkspaceRole) error {
	isSelfLeave := targetUserID == callerUserID

	if !isSelfLeave && !callerRole.CanManageMembers() {
		return domain.ErrForbidden
	}

	target, err := s.workspaceRepo.GetMember(ctx, workspaceID, targetUserID)
	if err != nil {
		return err
	}
	if target.Role == domain.RoleOwner && !isSelfLeave {
		return domain.ErrForbidden
	}
	if !isSelfLeave && !callerRole.IsSuperiorTo(target.Role) {
		return domain.ErrForbidden
	}

	// Cascade: remove from all projects in this workspace first
	if s.projectMemberRepo != nil {
		_ = s.projectMemberRepo.RemoveByWorkspace(ctx, workspaceID, targetUserID)
	}

	return s.workspaceRepo.RemoveMember(ctx, workspaceID, targetUserID)
}

func (s *WorkspaceService) GetStats(ctx context.Context, workspaceID uuid.UUID) (*dto.WorkspaceStatsResponse, error) {
	projectCount, err := s.workspaceRepo.CountProjects(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	byStatus, err := s.taskRepo.CountByStatusForWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	byAssigneeRaw, err := s.taskRepo.CountByAssigneeForWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	overdueCount, err := s.taskRepo.CountOverdueForWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	byAssignee := make([]dto.AssigneeStats, len(byAssigneeRaw))
	for i, a := range byAssigneeRaw {
		byAssignee[i] = dto.AssigneeStats{
			UserID:   a.UserID.String(),
			UserName: a.UserName,
			Count:    a.Count,
		}
	}

	return &dto.WorkspaceStatsResponse{
		ProjectCount:    projectCount,
		TasksByStatus:   byStatus,
		TasksByAssignee: byAssignee,
		OverdueCount:    overdueCount,
	}, nil
}

func (s *WorkspaceService) GetMember(ctx context.Context, workspaceID, userID uuid.UUID) (*domain.WorkspaceMember, error) {
	return s.workspaceRepo.GetMember(ctx, workspaceID, userID)
}

func (s *WorkspaceService) DirectAddMember(ctx context.Context, workspaceID, userID uuid.UUID, role domain.WorkspaceRole) error {
	// Check if already a member
	_, err := s.workspaceRepo.GetMember(ctx, workspaceID, userID)
	if err == nil {
		return domain.ErrConflict
	}

	now := time.Now().UTC()
	member := &domain.WorkspaceMember{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        role,
		JoinedAt:    now,
	}
	return s.workspaceRepo.AddMember(ctx, member)
}
