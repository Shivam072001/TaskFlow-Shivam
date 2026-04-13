package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type OrgInvitationService struct {
	invRepo *repository.OrgInvitationRepository
	orgRepo *repository.OrganizationRepository
	userRepo *repository.UserRepository
}

func NewOrgInvitationService(
	ir *repository.OrgInvitationRepository,
	or *repository.OrganizationRepository,
	ur *repository.UserRepository,
) *OrgInvitationService {
	return &OrgInvitationService{invRepo: ir, orgRepo: or, userRepo: ur}
}

func (s *OrgInvitationService) SendInvite(ctx context.Context, orgID, inviterID uuid.UUID, email string, role domain.OrgRole) (*domain.OrgInvitation, error) {
	// Check if user is already a member
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		_, err := s.orgRepo.GetMember(ctx, orgID, existingUser.ID)
		if err == nil {
			return nil, domain.ErrConflict
		}
	}

	var inviteeID *uuid.UUID
	if existingUser != nil {
		inviteeID = &existingUser.ID
	}

	now := time.Now().UTC()
	inv := &domain.OrgInvitation{
		ID:           uuid.New(),
		OrgID:        orgID,
		InviterID:    inviterID,
		InviteeEmail: email,
		InviteeID:    inviteeID,
		Role:         role,
		Status:       "pending",
		CreatedAt:    now,
	}

	if err := s.invRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("creating org invitation: %w", err)
	}
	return inv, nil
}

func (s *OrgInvitationService) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.OrgInvitation, error) {
	return s.invRepo.ListByOrg(ctx, orgID)
}

func (s *OrgInvitationService) ListByEmail(ctx context.Context, email string) ([]domain.OrgInvitation, error) {
	return s.invRepo.ListByEmail(ctx, email)
}

func (s *OrgInvitationService) Respond(ctx context.Context, invID, userID uuid.UUID, accept bool) error {
	inv, err := s.invRepo.GetByID(ctx, invID)
	if err != nil {
		return err
	}

	if inv.Status != "pending" {
		return domain.ErrConflict
	}

	status := "declined"
	if accept {
		status = "accepted"
	}

	if err := s.invRepo.Respond(ctx, invID, status, userID); err != nil {
		return fmt.Errorf("responding to invitation: %w", err)
	}

	if accept {
		now := time.Now().UTC()
		member := &domain.OrgMember{
			ID:       uuid.New(),
			OrgID:    inv.OrgID,
			UserID:   userID,
			Role:     inv.Role,
			JoinedAt: now,
		}
		if err := s.orgRepo.AddMember(ctx, member); err != nil {
			return fmt.Errorf("adding org member: %w", err)
		}
	}

	return nil
}
