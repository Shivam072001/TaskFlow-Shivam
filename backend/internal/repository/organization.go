package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type OrganizationRepository struct {
	pool *pgxpool.Pool
}

func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{pool: pool}
}

func (r *OrganizationRepository) Create(ctx context.Context, o *domain.Organization) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO organizations (id, name, slug, created_by, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		o.ID, o.Name, o.Slug, o.CreatedBy, o.CreatedAt,
	)
	return err
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	o := &domain.Organization{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, created_by, created_at FROM organizations WHERE id = $1`, id,
	).Scan(&o.ID, &o.Name, &o.Slug, &o.CreatedBy, &o.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return o, err
}

func (r *OrganizationRepository) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	o := &domain.Organization{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, slug, created_by, created_at FROM organizations WHERE slug = $1`, slug,
	).Scan(&o.ID, &o.Name, &o.Slug, &o.CreatedBy, &o.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return o, err
}

func (r *OrganizationRepository) Update(ctx context.Context, o *domain.Organization) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE organizations SET name = $2 WHERE id = $1`, o.ID, o.Name,
	)
	return err
}

func (r *OrganizationRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.Organization, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT o.id, o.name, o.slug, o.created_by, o.created_at
		 FROM organizations o
		 JOIN organization_members om ON om.org_id = o.id
		 WHERE om.user_id = $1
		 ORDER BY o.created_at ASC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []domain.Organization
	for rows.Next() {
		var o domain.Organization
		if err := rows.Scan(&o.ID, &o.Name, &o.Slug, &o.CreatedBy, &o.CreatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}
	return orgs, rows.Err()
}

func (r *OrganizationRepository) AddMember(ctx context.Context, m *domain.OrgMember) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO organization_members (id, org_id, user_id, role, joined_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		m.ID, m.OrgID, m.UserID, m.Role, m.JoinedAt,
	)
	return err
}

func (r *OrganizationRepository) GetMember(ctx context.Context, orgID, userID uuid.UUID) (*domain.OrgMember, error) {
	m := &domain.OrgMember{}
	err := r.pool.QueryRow(ctx,
		`SELECT om.id, om.org_id, om.user_id, om.role, om.joined_at, u.name, u.email
		 FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 WHERE om.org_id = $1 AND om.user_id = $2`,
		orgID, userID,
	).Scan(&m.ID, &m.OrgID, &m.UserID, &m.Role, &m.JoinedAt, &m.UserName, &m.UserEmail)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return m, err
}

func (r *OrganizationRepository) ListMembers(ctx context.Context, orgID uuid.UUID) ([]domain.OrgMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT om.id, om.org_id, om.user_id, om.role, om.joined_at, u.name, u.email
		 FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 WHERE om.org_id = $1
		 ORDER BY om.joined_at ASC`, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.OrgMember
	for rows.Next() {
		var m domain.OrgMember
		if err := rows.Scan(&m.ID, &m.OrgID, &m.UserID, &m.Role, &m.JoinedAt, &m.UserName, &m.UserEmail); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *OrganizationRepository) ListPrefixes(ctx context.Context, orgID uuid.UUID) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT prefix FROM projects WHERE org_id = $1 AND prefix != '' ORDER BY prefix`, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prefixes []string
	for rows.Next() {
		var p string
		if err := rows.Scan(&p); err != nil {
			return nil, err
		}
		prefixes = append(prefixes, p)
	}
	return prefixes, rows.Err()
}

func (r *OrganizationRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM organizations WHERE slug = $1)`, slug,
	).Scan(&exists)
	return exists, err
}

func (r *OrganizationRepository) UpdateMemberRole(ctx context.Context, orgID, userID uuid.UUID, role domain.OrgRole) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE organization_members SET role = $3 WHERE org_id = $1 AND user_id = $2`,
		orgID, userID, role,
	)
	return err
}

func (r *OrganizationRepository) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM organization_members WHERE org_id = $1 AND user_id = $2`, orgID, userID,
	)
	return err
}

func (r *OrganizationRepository) RemoveFromTeamsByOrg(ctx context.Context, orgID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM team_members
		 WHERE user_id = $2 AND team_id IN (SELECT id FROM teams WHERE org_id = $1)`,
		orgID, userID,
	)
	return err
}

func (r *OrganizationRepository) RemoveFromWorkspacesByOrg(ctx context.Context, orgID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM workspace_members
		 WHERE user_id = $2 AND workspace_id IN (SELECT id FROM workspaces WHERE org_id = $1)`,
		orgID, userID,
	)
	return err
}

func (r *OrganizationRepository) GetTaskByKey(ctx context.Context, orgID uuid.UUID, taskKey string) (*domain.Task, uuid.UUID, error) {
	t := &domain.Task{}
	var workspaceID uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT t.id, t.title, t.description, t.status, t.priority, t.project_id,
		        t.assignee_id, t.start_date, t.due_date, t.created_by, t.created_at, t.updated_at,
		        t.task_number, t.task_key, t.blocked_reason, t.blocked_by_task,
		        p.workspace_id
		 FROM tasks t
		 JOIN projects p ON p.id = t.project_id
		 WHERE p.org_id = $1 AND t.task_key = $2`,
		orgID, taskKey,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID,
		&t.AssigneeID, &t.StartDate, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		&t.TaskNumber, &t.TaskKey, &t.BlockedReason, &t.BlockedByTask,
		&workspaceID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, uuid.Nil, domain.ErrNotFound
	}
	return t, workspaceID, err
}
