package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type OrganizationService struct {
	orgRepo           *repository.OrganizationRepository
	projectMemberRepo *repository.ProjectMemberRepository
}

func NewOrganizationService(or *repository.OrganizationRepository, pmr *repository.ProjectMemberRepository) *OrganizationService {
	return &OrganizationService{orgRepo: or, projectMemberRepo: pmr}
}

func (s *OrganizationService) Create(ctx context.Context, userID uuid.UUID, name, slug string) (*domain.Organization, error) {
	taken, err := s.orgRepo.SlugExists(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("checking slug: %w", err)
	}
	if taken {
		return nil, domain.ErrConflict
	}

	now := time.Now().UTC()
	o := &domain.Organization{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		CreatedBy: userID,
		CreatedAt: now,
	}

	if err := s.orgRepo.Create(ctx, o); err != nil {
		return nil, fmt.Errorf("creating organization: %w", err)
	}

	member := &domain.OrgMember{
		ID:       uuid.New(),
		OrgID:    o.ID,
		UserID:   userID,
		Role:     domain.OrgRoleOwner,
		JoinedAt: now,
	}
	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("adding owner as org member: %w", err)
	}

	return o, nil
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	return s.orgRepo.GetByID(ctx, id)
}

func (s *OrganizationService) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	return s.orgRepo.GetBySlug(ctx, slug)
}

func (s *OrganizationService) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Organization, error) {
	return s.orgRepo.ListByUser(ctx, userID)
}

func (s *OrganizationService) Update(ctx context.Context, id uuid.UUID, role domain.OrgRole, name *string) (*domain.Organization, error) {
	if !role.CanManageOrg() {
		return nil, domain.ErrForbidden
	}

	o, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		o.Name = *name
	}

	if err := s.orgRepo.Update(ctx, o); err != nil {
		return nil, fmt.Errorf("updating organization: %w", err)
	}
	return o, nil
}

func (s *OrganizationService) GetMember(ctx context.Context, orgID, userID uuid.UUID) (*domain.OrgMember, error) {
	return s.orgRepo.GetMember(ctx, orgID, userID)
}

func (s *OrganizationService) ListMembers(ctx context.Context, orgID uuid.UUID) ([]domain.OrgMember, error) {
	return s.orgRepo.ListMembers(ctx, orgID)
}

func (s *OrganizationService) ListPrefixes(ctx context.Context, orgID uuid.UUID) ([]string, error) {
	return s.orgRepo.ListPrefixes(ctx, orgID)
}

func (s *OrganizationService) GetTaskByKey(ctx context.Context, orgID uuid.UUID, taskKey string) (*domain.Task, uuid.UUID, error) {
	return s.orgRepo.GetTaskByKey(ctx, orgID, taskKey)
}

// RemoveMemberCascade removes a user from the org and all descendant entities:
// project_members -> workspace_members -> team_members -> organization_members
func (s *OrganizationService) RemoveMemberCascade(ctx context.Context, orgID, userID uuid.UUID) error {
	if s.projectMemberRepo != nil {
		_ = s.projectMemberRepo.RemoveByOrg(ctx, orgID, userID)
	}
	_ = s.orgRepo.RemoveFromWorkspacesByOrg(ctx, orgID, userID)
	_ = s.orgRepo.RemoveFromTeamsByOrg(ctx, orgID, userID)
	return s.orgRepo.RemoveMember(ctx, orgID, userID)
}
