package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type OrgInvitationRepository struct {
	pool *pgxpool.Pool
}

func NewOrgInvitationRepository(pool *pgxpool.Pool) *OrgInvitationRepository {
	return &OrgInvitationRepository{pool: pool}
}

func (r *OrgInvitationRepository) Create(ctx context.Context, inv *domain.OrgInvitation) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO org_invitations (id, org_id, inviter_id, invitee_email, invitee_id, role, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		inv.ID, inv.OrgID, inv.InviterID, inv.InviteeEmail, inv.InviteeID, inv.Role, inv.Status, inv.CreatedAt,
	)
	return err
}

func (r *OrgInvitationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.OrgInvitation, error) {
	inv := &domain.OrgInvitation{}
	err := r.pool.QueryRow(ctx,
		`SELECT oi.id, oi.org_id, oi.inviter_id, oi.invitee_email, oi.invitee_id,
		        oi.role, oi.status, oi.created_at, oi.responded_at,
		        o.name, u.name
		 FROM org_invitations oi
		 JOIN organizations o ON o.id = oi.org_id
		 JOIN users u ON u.id = oi.inviter_id
		 WHERE oi.id = $1`, id,
	).Scan(
		&inv.ID, &inv.OrgID, &inv.InviterID, &inv.InviteeEmail, &inv.InviteeID,
		&inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt,
		&inv.OrgName, &inv.InviterName,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return inv, err
}

func (r *OrgInvitationRepository) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.OrgInvitation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT oi.id, oi.org_id, oi.inviter_id, oi.invitee_email, oi.invitee_id,
		        oi.role, oi.status, oi.created_at, oi.responded_at,
		        o.name, u.name
		 FROM org_invitations oi
		 JOIN organizations o ON o.id = oi.org_id
		 JOIN users u ON u.id = oi.inviter_id
		 WHERE oi.org_id = $1
		 ORDER BY oi.created_at DESC`, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []domain.OrgInvitation
	for rows.Next() {
		var inv domain.OrgInvitation
		if err := rows.Scan(
			&inv.ID, &inv.OrgID, &inv.InviterID, &inv.InviteeEmail, &inv.InviteeID,
			&inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt,
			&inv.OrgName, &inv.InviterName,
		); err != nil {
			return nil, err
		}
		invitations = append(invitations, inv)
	}
	return invitations, rows.Err()
}

func (r *OrgInvitationRepository) ListByEmail(ctx context.Context, email string) ([]domain.OrgInvitation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT oi.id, oi.org_id, oi.inviter_id, oi.invitee_email, oi.invitee_id,
		        oi.role, oi.status, oi.created_at, oi.responded_at,
		        o.name, u.name
		 FROM org_invitations oi
		 JOIN organizations o ON o.id = oi.org_id
		 JOIN users u ON u.id = oi.inviter_id
		 WHERE oi.invitee_email = $1 AND oi.status = 'pending'
		 ORDER BY oi.created_at DESC`, email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []domain.OrgInvitation
	for rows.Next() {
		var inv domain.OrgInvitation
		if err := rows.Scan(
			&inv.ID, &inv.OrgID, &inv.InviterID, &inv.InviteeEmail, &inv.InviteeID,
			&inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt,
			&inv.OrgName, &inv.InviterName,
		); err != nil {
			return nil, err
		}
		invitations = append(invitations, inv)
	}
	return invitations, rows.Err()
}

func (r *OrgInvitationRepository) Respond(ctx context.Context, id uuid.UUID, status string, inviteeID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE org_invitations SET status = $2, invitee_id = $3, responded_at = now() WHERE id = $1`,
		id, status, inviteeID,
	)
	return err
}
