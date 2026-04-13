package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type InvitationRepository struct {
	pool *pgxpool.Pool
}

func NewInvitationRepository(pool *pgxpool.Pool) *InvitationRepository {
	return &InvitationRepository{pool: pool}
}

func (r *InvitationRepository) Create(ctx context.Context, inv *domain.WorkspaceInvitation) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO workspace_invitations (id, workspace_id, inviter_id, invitee_email, invitee_id, role, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		inv.ID, inv.WorkspaceID, inv.InviterID, inv.InviteeEmail,
		inv.InviteeID, inv.Role, inv.Status, inv.CreatedAt,
	)
	return err
}

func (r *InvitationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WorkspaceInvitation, error) {
	inv := &domain.WorkspaceInvitation{}
	err := r.pool.QueryRow(ctx,
		`SELECT wi.id, wi.workspace_id, wi.inviter_id, wi.invitee_email, wi.invitee_id,
		        wi.role, wi.status, wi.created_at, wi.responded_at,
		        w.name, u.name
		 FROM workspace_invitations wi
		 JOIN workspaces w ON w.id = wi.workspace_id
		 JOIN users u ON u.id = wi.inviter_id
		 WHERE wi.id = $1`,
		id,
	).Scan(&inv.ID, &inv.WorkspaceID, &inv.InviterID, &inv.InviteeEmail, &inv.InviteeID,
		&inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt,
		&inv.WorkspaceName, &inv.InviterName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return inv, err
}

func (r *InvitationRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.WorkspaceInvitation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT wi.id, wi.workspace_id, wi.inviter_id, wi.invitee_email, wi.invitee_id,
		        wi.role, wi.status, wi.created_at, wi.responded_at,
		        w.name, u.name
		 FROM workspace_invitations wi
		 JOIN workspaces w ON w.id = wi.workspace_id
		 JOIN users u ON u.id = wi.inviter_id
		 WHERE wi.workspace_id = $1
		 ORDER BY wi.created_at DESC`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []domain.WorkspaceInvitation
	for rows.Next() {
		var inv domain.WorkspaceInvitation
		if err := rows.Scan(&inv.ID, &inv.WorkspaceID, &inv.InviterID, &inv.InviteeEmail, &inv.InviteeID,
			&inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt,
			&inv.WorkspaceName, &inv.InviterName); err != nil {
			return nil, err
		}
		invitations = append(invitations, inv)
	}
	return invitations, rows.Err()
}

func (r *InvitationRepository) ListPendingByUser(ctx context.Context, email string) ([]domain.WorkspaceInvitation, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT wi.id, wi.workspace_id, wi.inviter_id, wi.invitee_email, wi.invitee_id,
		        wi.role, wi.status, wi.created_at, wi.responded_at,
		        w.name, u.name
		 FROM workspace_invitations wi
		 JOIN workspaces w ON w.id = wi.workspace_id
		 JOIN users u ON u.id = wi.inviter_id
		 WHERE wi.invitee_email = $1 AND wi.status = 'pending'
		 ORDER BY wi.created_at DESC`,
		email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []domain.WorkspaceInvitation
	for rows.Next() {
		var inv domain.WorkspaceInvitation
		if err := rows.Scan(&inv.ID, &inv.WorkspaceID, &inv.InviterID, &inv.InviteeEmail, &inv.InviteeID,
			&inv.Role, &inv.Status, &inv.CreatedAt, &inv.RespondedAt,
			&inv.WorkspaceName, &inv.InviterName); err != nil {
			return nil, err
		}
		invitations = append(invitations, inv)
	}
	return invitations, rows.Err()
}

func (r *InvitationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.InvitationStatus, respondedAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE workspace_invitations SET status = $2, responded_at = $3 WHERE id = $1`,
		id, status, respondedAt,
	)
	return err
}

func (r *InvitationRepository) HasPendingInvite(ctx context.Context, workspaceID uuid.UUID, email string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM workspace_invitations
			WHERE workspace_id = $1 AND invitee_email = $2 AND status = 'pending'
		)`,
		workspaceID, email,
	).Scan(&exists)
	return exists, err
}
