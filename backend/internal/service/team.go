package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type TeamService struct {
	teamRepo      *repository.TeamRepository
	workspaceRepo *repository.WorkspaceRepository
	projectMemberRepo *repository.ProjectMemberRepository
	orgRepo       *repository.OrganizationRepository
}

func NewTeamService(
	tr *repository.TeamRepository,
	wr *repository.WorkspaceRepository,
	pmr *repository.ProjectMemberRepository,
	or *repository.OrganizationRepository,
) *TeamService {
	return &TeamService{teamRepo: tr, workspaceRepo: wr, projectMemberRepo: pmr, orgRepo: or}
}

func (s *TeamService) Create(ctx context.Context, orgID, userID uuid.UUID, name string) (*domain.Team, error) {
	now := time.Now().UTC()
	t := &domain.Team{
		ID:        uuid.New(),
		OrgID:     orgID,
		Name:      name,
		CreatedBy: userID,
		CreatedAt: now,
	}
	if err := s.teamRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("creating team: %w", err)
	}
	return t, nil
}

func (s *TeamService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	return s.teamRepo.GetByID(ctx, id)
}

func (s *TeamService) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.Team, error) {
	return s.teamRepo.ListByOrg(ctx, orgID)
}

func (s *TeamService) Update(ctx context.Context, teamID uuid.UUID, name string) error {
	return s.teamRepo.Update(ctx, teamID, name)
}

func (s *TeamService) Delete(ctx context.Context, teamID uuid.UUID) error {
	return s.teamRepo.Delete(ctx, teamID)
}

func (s *TeamService) ListMembers(ctx context.Context, teamID uuid.UUID) ([]domain.TeamMember, error) {
	return s.teamRepo.ListMembers(ctx, teamID)
}

// AddMember adds a user to a team and auto-cascades into workspace_members and project_members
// for any workspaces/projects the team is already assigned to.
func (s *TeamService) AddMember(ctx context.Context, teamID, userID uuid.UUID) (*domain.TeamMember, error) {
	now := time.Now().UTC()
	m := &domain.TeamMember{
		ID:      uuid.New(),
		TeamID:  teamID,
		UserID:  userID,
		AddedAt: now,
	}
	if err := s.teamRepo.AddMember(ctx, m); err != nil {
		return nil, fmt.Errorf("adding team member: %w", err)
	}

	// Auto-cascade: add to all workspaces this team is assigned to
	wsAssignments, err := s.teamRepo.ListWorkspaceTeamAssignments(ctx, teamID)
	if err != nil {
		return m, nil // member added, cascade is best-effort
	}
	for _, wt := range wsAssignments {
		_ = s.workspaceRepo.AddMember(ctx, &domain.WorkspaceMember{
			ID:          uuid.New(),
			WorkspaceID: wt.WorkspaceID,
			UserID:      userID,
			Role:        domain.WorkspaceRole(wt.DefaultRole),
			JoinedAt:    now,
		})
	}

	// Auto-cascade: add to all projects this team is assigned to
	ptAssignments, err := s.teamRepo.ListProjectTeamAssignments(ctx, teamID)
	if err != nil {
		return m, nil
	}
	for _, pt := range ptAssignments {
		_ = s.projectMemberRepo.Add(ctx, &domain.ProjectMember{
			ID:        uuid.New(),
			ProjectID: pt.ProjectID,
			UserID:    userID,
			Role:      domain.WorkspaceRole(pt.DefaultRole),
			JoinedAt:  now,
		})
	}

	return m, nil
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	return s.teamRepo.RemoveMember(ctx, teamID, userID)
}

// AddTeamToWorkspace assigns a team to a workspace and adds all team members as workspace members.
func (s *TeamService) AddTeamToWorkspace(ctx context.Context, workspaceID, teamID uuid.UUID, defaultRole string) error {
	now := time.Now().UTC()
	wt := &domain.WorkspaceTeam{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		TeamID:      teamID,
		DefaultRole: defaultRole,
		AddedAt:     now,
	}
	if err := s.teamRepo.AddTeamToWorkspace(ctx, wt); err != nil {
		return fmt.Errorf("adding team to workspace: %w", err)
	}

	members, err := s.teamRepo.ListMembers(ctx, teamID)
	if err != nil {
		return nil
	}
	for _, m := range members {
		_ = s.workspaceRepo.AddMember(ctx, &domain.WorkspaceMember{
			ID:          uuid.New(),
			WorkspaceID: workspaceID,
			UserID:      m.UserID,
			Role:        domain.WorkspaceRole(defaultRole),
			JoinedAt:    now,
		})
	}
	return nil
}

func (s *TeamService) RemoveTeamFromWorkspace(ctx context.Context, workspaceID, teamID uuid.UUID) error {
	return s.teamRepo.RemoveTeamFromWorkspace(ctx, workspaceID, teamID)
}

// AddTeamToProject assigns a team to a project and adds all team members (who are workspace members) as project members.
func (s *TeamService) AddTeamToProject(ctx context.Context, projectID, workspaceID, teamID uuid.UUID, defaultRole string) error {
	now := time.Now().UTC()
	pt := &domain.ProjectTeam{
		ID:          uuid.New(),
		ProjectID:   projectID,
		TeamID:      teamID,
		DefaultRole: defaultRole,
		AddedAt:     now,
	}
	if err := s.teamRepo.AddTeamToProject(ctx, pt); err != nil {
		return fmt.Errorf("adding team to project: %w", err)
	}

	members, err := s.teamRepo.ListMembers(ctx, teamID)
	if err != nil {
		return nil
	}
	for _, m := range members {
		// Only add if user is a workspace member
		_, wsErr := s.workspaceRepo.GetMember(ctx, workspaceID, m.UserID)
		if wsErr != nil {
			continue
		}
		_ = s.projectMemberRepo.Add(ctx, &domain.ProjectMember{
			ID:        uuid.New(),
			ProjectID: projectID,
			UserID:    m.UserID,
			Role:      domain.WorkspaceRole(defaultRole),
			JoinedAt:  now,
		})
	}
	return nil
}

func (s *TeamService) RemoveTeamFromProject(ctx context.Context, projectID, teamID uuid.UUID) error {
	return s.teamRepo.RemoveTeamFromProject(ctx, projectID, teamID)
}

func (s *TeamService) ListTeamsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.Team, error) {
	return s.teamRepo.ListTeamsByWorkspace(ctx, workspaceID)
}

func (s *TeamService) ListTeamsByProject(ctx context.Context, projectID uuid.UUID) ([]domain.Team, error) {
	return s.teamRepo.ListTeamsByProject(ctx, projectID)
}
