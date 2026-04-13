package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type WorkspaceRepository struct {
	pool *pgxpool.Pool
}

func NewWorkspaceRepository(pool *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{pool: pool}
}

func (r *WorkspaceRepository) Create(ctx context.Context, w *domain.Workspace) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO workspaces (id, org_id, name, description, created_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		w.ID, w.OrgID, w.Name, w.Description, w.CreatedBy, w.CreatedAt, w.UpdatedAt,
	)
	return err
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Workspace, error) {
	w := &domain.Workspace{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, org_id, name, description, created_by, created_at, updated_at
		 FROM workspaces WHERE id = $1`,
		id,
	).Scan(&w.ID, &w.OrgID, &w.Name, &w.Description, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return w, err
}

func (r *WorkspaceRepository) ListByUser(ctx context.Context, userID, orgID uuid.UUID, page, limit int) ([]domain.Workspace, int, error) {
	baseWhere := `FROM workspaces w
		 JOIN workspace_members wm ON wm.workspace_id = w.id
		 WHERE wm.user_id = $1 AND w.org_id = $2`

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) "+baseWhere, userID, orgID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	rows, err := r.pool.Query(ctx,
		`SELECT w.id, w.org_id, w.name, w.description, w.created_by, w.created_at, w.updated_at `+
			baseWhere+` ORDER BY w.created_at DESC LIMIT $3 OFFSET $4`,
		userID, orgID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var workspaces []domain.Workspace
	for rows.Next() {
		var w domain.Workspace
		if err := rows.Scan(&w.ID, &w.OrgID, &w.Name, &w.Description, &w.CreatedBy, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, 0, err
		}
		workspaces = append(workspaces, w)
	}
	return workspaces, total, rows.Err()
}

func (r *WorkspaceRepository) Update(ctx context.Context, w *domain.Workspace) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE workspaces SET name = $2, description = $3, updated_at = $4 WHERE id = $1`,
		w.ID, w.Name, w.Description, w.UpdatedAt,
	)
	return err
}

func (r *WorkspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM workspaces WHERE id = $1`, id)
	return err
}

func (r *WorkspaceRepository) AddMember(ctx context.Context, m *domain.WorkspaceMember) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO workspace_members (id, workspace_id, user_id, role, joined_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (workspace_id, user_id) DO NOTHING`,
		m.ID, m.WorkspaceID, m.UserID, m.Role, m.JoinedAt,
	)
	return err
}

func (r *WorkspaceRepository) GetMember(ctx context.Context, workspaceID, userID uuid.UUID) (*domain.WorkspaceMember, error) {
	m := &domain.WorkspaceMember{}
	err := r.pool.QueryRow(ctx,
		`SELECT wm.id, wm.workspace_id, wm.user_id, wm.role, wm.joined_at, u.name, u.email
		 FROM workspace_members wm
		 JOIN users u ON u.id = wm.user_id
		 WHERE wm.workspace_id = $1 AND wm.user_id = $2`,
		workspaceID, userID,
	).Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.JoinedAt, &m.UserName, &m.UserEmail)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return m, err
}

func (r *WorkspaceRepository) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]domain.WorkspaceMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT wm.id, wm.workspace_id, wm.user_id, wm.role, wm.joined_at, u.name, u.email
		 FROM workspace_members wm
		 JOIN users u ON u.id = wm.user_id
		 WHERE wm.workspace_id = $1
		 ORDER BY wm.joined_at ASC`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.WorkspaceMember
	for rows.Next() {
		var m domain.WorkspaceMember
		if err := rows.Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.JoinedAt, &m.UserName, &m.UserEmail); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *WorkspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role domain.WorkspaceRole) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE workspace_members SET role = $3 WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID, userID, role,
	)
	return err
}

func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID, userID,
	)
	return err
}

func (r *WorkspaceRepository) CountProjects(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM projects WHERE workspace_id = $1`, workspaceID,
	).Scan(&count)
	return count, err
}

func (r *WorkspaceRepository) CountMembersByWorkspace(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM workspace_members WHERE workspace_id = $1`, workspaceID,
	).Scan(&count)
	return count, err
}
