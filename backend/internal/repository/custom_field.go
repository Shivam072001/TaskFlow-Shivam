package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type CustomFieldRepository struct {
	pool *pgxpool.Pool
}

func NewCustomFieldRepository(pool *pgxpool.Pool) *CustomFieldRepository {
	return &CustomFieldRepository{pool: pool}
}

func (r *CustomFieldRepository) CreateDefinition(ctx context.Context, def *domain.CustomFieldDefinition) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO custom_field_definitions (id, project_id, name, field_type, options, required, created_by, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		def.ID, def.ProjectID, def.Name, def.FieldType, def.Options, def.Required, def.CreatedBy, def.CreatedAt,
	)
	return err
}

func (r *CustomFieldRepository) GetDefinitionByID(ctx context.Context, id uuid.UUID) (*domain.CustomFieldDefinition, error) {
	def := &domain.CustomFieldDefinition{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, project_id, name, field_type, options, required, created_by, created_at
		 FROM custom_field_definitions WHERE id = $1`,
		id,
	).Scan(&def.ID, &def.ProjectID, &def.Name, &def.FieldType, &def.Options, &def.Required, &def.CreatedBy, &def.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return def, err
}

func (r *CustomFieldRepository) ListDefinitions(ctx context.Context, projectID uuid.UUID) ([]domain.CustomFieldDefinition, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, name, field_type, options, required, created_by, created_at
		 FROM custom_field_definitions WHERE project_id = $1
		 ORDER BY created_at ASC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var defs []domain.CustomFieldDefinition
	for rows.Next() {
		var d domain.CustomFieldDefinition
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.Name, &d.FieldType, &d.Options, &d.Required, &d.CreatedBy, &d.CreatedAt); err != nil {
			return nil, err
		}
		defs = append(defs, d)
	}
	return defs, rows.Err()
}

func (r *CustomFieldRepository) DeleteDefinition(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM custom_field_definitions WHERE id = $1`, id)
	return err
}

func (r *CustomFieldRepository) UpsertValue(ctx context.Context, val *domain.CustomFieldValue) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO custom_field_values (id, task_id, field_id, value)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (task_id, field_id) DO UPDATE SET value = EXCLUDED.value`,
		val.ID, val.TaskID, val.FieldID, val.Value,
	)
	return err
}

func (r *CustomFieldRepository) GetValues(ctx context.Context, taskID uuid.UUID) ([]domain.CustomFieldValue, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, task_id, field_id, value
		 FROM custom_field_values WHERE task_id = $1`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vals []domain.CustomFieldValue
	for rows.Next() {
		var v domain.CustomFieldValue
		if err := rows.Scan(&v.ID, &v.TaskID, &v.FieldID, &v.Value); err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, rows.Err()
}

func (r *CustomFieldRepository) GetWIPLimits(ctx context.Context, projectID uuid.UUID) ([]domain.WIPLimit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, status, max_tasks
		 FROM project_wip_limits WHERE project_id = $1
		 ORDER BY status ASC`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var limits []domain.WIPLimit
	for rows.Next() {
		var l domain.WIPLimit
		if err := rows.Scan(&l.ID, &l.ProjectID, &l.Status, &l.MaxTasks); err != nil {
			return nil, err
		}
		limits = append(limits, l)
	}
	return limits, rows.Err()
}

func (r *CustomFieldRepository) UpsertWIPLimit(ctx context.Context, limit *domain.WIPLimit) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO project_wip_limits (id, project_id, status, max_tasks)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (project_id, status) DO UPDATE SET max_tasks = EXCLUDED.max_tasks`,
		limit.ID, limit.ProjectID, limit.Status, limit.MaxTasks,
	)
	return err
}

func (r *CustomFieldRepository) DeleteWIPLimit(ctx context.Context, projectID uuid.UUID, status domain.TaskStatus) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM project_wip_limits WHERE project_id = $1 AND status = $2`,
		projectID, status,
	)
	return err
}

func (r *CustomFieldRepository) CountTasksByStatus(ctx context.Context, projectID uuid.UUID, status domain.TaskStatus) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks WHERE project_id = $1 AND status = $2`,
		projectID, status,
	).Scan(&count)
	return count, err
}
