package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type ProjectFilter struct {
	WorkspaceID uuid.UUID
	Search      string
	OwnerID     string
	Page        int
	Limit       int
}

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) Create(ctx context.Context, p *domain.Project) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO projects (id, org_id, name, prefix, description, workspace_id, owner_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		p.ID, p.OrgID, p.Name, p.Prefix, p.Description, p.WorkspaceID, p.OwnerID, p.CreatedAt,
	)
	return err
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	p := &domain.Project{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, org_id, name, prefix, description, workspace_id, owner_id, created_at
		 FROM projects WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.OrgID, &p.Name, &p.Prefix, &p.Description, &p.WorkspaceID, &p.OwnerID, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *ProjectRepository) ListByWorkspace(ctx context.Context, f ProjectFilter) ([]domain.Project, int, error) {
	where := "WHERE p.workspace_id = $1"
	args := []interface{}{f.WorkspaceID}
	argIdx := 2

	if f.Search != "" {
		where += fmt.Sprintf(" AND (p.name ILIKE '%%' || $%d || '%%' OR p.description ILIKE '%%' || $%d || '%%')", argIdx, argIdx)
		args = append(args, f.Search)
		argIdx++
	}
	if f.OwnerID != "" {
		uid, err := uuid.Parse(f.OwnerID)
		if err == nil {
			where += fmt.Sprintf(" AND p.owner_id = $%d", argIdx)
			args = append(args, uid)
			argIdx++
		}
	}

	var total int
	if err := r.pool.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM projects p %s", where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	query := fmt.Sprintf(
		`SELECT p.id, p.org_id, p.name, p.prefix, p.description, p.workspace_id, p.owner_id, p.created_at
		 FROM projects p %s
		 ORDER BY p.created_at DESC
		 LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	args = append(args, f.Limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(&p.ID, &p.OrgID, &p.Name, &p.Prefix, &p.Description, &p.WorkspaceID, &p.OwnerID, &p.CreatedAt); err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}
	return projects, total, rows.Err()
}

func (r *ProjectRepository) Update(ctx context.Context, p *domain.Project) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE projects SET name = $2, description = $3 WHERE id = $1`,
		p.ID, p.Name, p.Description,
	)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id = $1`, id)
	return err
}
