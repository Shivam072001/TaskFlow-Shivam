package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type ProjectMemberRepository struct {
	pool *pgxpool.Pool
}

func NewProjectMemberRepository(pool *pgxpool.Pool) *ProjectMemberRepository {
	return &ProjectMemberRepository{pool: pool}
}

func (r *ProjectMemberRepository) Add(ctx context.Context, m *domain.ProjectMember) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO project_members (id, project_id, user_id, role, joined_at) VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (project_id, user_id) DO NOTHING`,
		m.ID, m.ProjectID, m.UserID, m.Role, m.JoinedAt,
	)
	return err
}

func (r *ProjectMemberRepository) Remove(ctx context.Context, projectID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`, projectID, userID,
	)
	return err
}

func (r *ProjectMemberRepository) GetMember(ctx context.Context, projectID, userID uuid.UUID) (*domain.ProjectMember, error) {
	m := &domain.ProjectMember{}
	err := r.pool.QueryRow(ctx,
		`SELECT pm.id, pm.project_id, pm.user_id, pm.role, pm.joined_at, u.name, u.email
		 FROM project_members pm
		 JOIN users u ON u.id = pm.user_id
		 WHERE pm.project_id = $1 AND pm.user_id = $2`,
		projectID, userID,
	).Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.JoinedAt, &m.UserName, &m.UserEmail)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return m, err
}

func (r *ProjectMemberRepository) List(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT pm.id, pm.project_id, pm.user_id, pm.role, pm.joined_at, u.name, u.email
		 FROM project_members pm
		 JOIN users u ON u.id = pm.user_id
		 WHERE pm.project_id = $1
		 ORDER BY pm.joined_at`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.ProjectMember
	for rows.Next() {
		var m domain.ProjectMember
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Role, &m.JoinedAt, &m.UserName, &m.UserEmail); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *ProjectMemberRepository) UpdateRole(ctx context.Context, projectID, userID uuid.UUID, role domain.WorkspaceRole) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE project_members SET role = $3 WHERE project_id = $1 AND user_id = $2`,
		projectID, userID, role,
	)
	return err
}

// RemoveByWorkspace removes all project members for a user across all projects in a workspace.
func (r *ProjectMemberRepository) RemoveByWorkspace(ctx context.Context, workspaceID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM project_members
		 WHERE user_id = $2 AND project_id IN (SELECT id FROM projects WHERE workspace_id = $1)`,
		workspaceID, userID,
	)
	return err
}

// RemoveByOrg removes all project members for a user across all projects in an org.
func (r *ProjectMemberRepository) RemoveByOrg(ctx context.Context, orgID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM project_members
		 WHERE user_id = $2 AND project_id IN (SELECT id FROM projects WHERE org_id = $1)`,
		orgID, userID,
	)
	return err
}
