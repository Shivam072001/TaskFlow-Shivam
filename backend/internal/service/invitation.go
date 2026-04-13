package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type InvitationService struct {
	invitationRepo *repository.InvitationRepository
	workspaceRepo  *repository.WorkspaceRepository
	userRepo       *repository.UserRepository
}

func NewInvitationService(ir *repository.InvitationRepository, wr *repository.WorkspaceRepository, ur *repository.UserRepository) *InvitationService {
	return &InvitationService{invitationRepo: ir, workspaceRepo: wr, userRepo: ur}
}

func (s *InvitationService) SendInvite(ctx context.Context, workspaceID, inviterID uuid.UUID, callerRole domain.WorkspaceRole, email string, role domain.WorkspaceRole) (*domain.WorkspaceInvitation, error) {
	if !callerRole.CanInviteMembers() {
		return nil, domain.ErrForbidden
	}
	if !callerRole.IsSuperiorTo(role) {
		return nil, fmt.Errorf("cannot assign role equal or above your own: %w", domain.ErrForbidden)
	}
	if !role.IsAssignable() {
		return nil, domain.ErrInvalidInput
	}

	hasPending, err := s.invitationRepo.HasPendingInvite(ctx, workspaceID, email)
	if err != nil {
		return nil, fmt.Errorf("checking pending invite: %w", err)
	}
	if hasPending {
		return nil, domain.ErrConflict
	}

	var inviteeID *uuid.UUID
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("looking up user: %w", err)
	}
	if user != nil {
		existing, err := s.workspaceRepo.GetMember(ctx, workspaceID, user.ID)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, domain.ErrConflict
		}
		inviteeID = &user.ID
	}

	inv := &domain.WorkspaceInvitation{
		ID:           uuid.New(),
		WorkspaceID:  workspaceID,
		InviterID:    inviterID,
		InviteeEmail: email,
		InviteeID:    inviteeID,
		Role:         role,
		Status:       domain.InvitePending,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.invitationRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("creating invitation: %w", err)
	}
	return inv, nil
}

func (s *InvitationService) ListPending(ctx context.Context, userEmail string) ([]domain.WorkspaceInvitation, error) {
	return s.invitationRepo.ListPendingByUser(ctx, userEmail)
}

func (s *InvitationService) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.WorkspaceInvitation, error) {
	return s.invitationRepo.ListByWorkspace(ctx, workspaceID)
}

func (s *InvitationService) Respond(ctx context.Context, invitationID, userID uuid.UUID, userEmail string, accept bool) error {
	inv, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return err
	}

	if inv.InviteeEmail != userEmail {
		return domain.ErrForbidden
	}
	if inv.Status != domain.InvitePending {
		return fmt.Errorf("invitation already %s: %w", inv.Status, domain.ErrConflict)
	}

	now := time.Now().UTC()
	if accept {
		member := &domain.WorkspaceMember{
			ID:          uuid.New(),
			WorkspaceID: inv.WorkspaceID,
			UserID:      userID,
			Role:        inv.Role,
			JoinedAt:    now,
		}
		if err := s.workspaceRepo.AddMember(ctx, member); err != nil {
			return fmt.Errorf("adding member: %w", err)
		}
		return s.invitationRepo.UpdateStatus(ctx, invitationID, domain.InviteAccepted, now)
	}

	return s.invitationRepo.UpdateStatus(ctx, invitationID, domain.InviteDeclined, now)
}
