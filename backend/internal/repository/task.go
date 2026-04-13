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

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) NextTaskNumber(ctx context.Context, projectID uuid.UUID) (int, error) {
	var num int
	err := r.pool.QueryRow(ctx,
		`SELECT COALESCE(MAX(task_number), 0) + 1 FROM tasks WHERE project_id = $1`,
		projectID,
	).Scan(&num)
	return num, err
}

func (r *TaskRepository) Create(ctx context.Context, t *domain.Task) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, start_date, due_date, created_by, created_at, updated_at, task_number, task_key, blocked_reason, blocked_by_task)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		t.ID, t.Title, t.Description, t.Status, t.Priority, t.ProjectID,
		t.AssigneeID, t.StartDate, t.DueDate, t.CreatedBy, t.CreatedAt, t.UpdatedAt,
		t.TaskNumber, t.TaskKey, t.BlockedReason, t.BlockedByTask,
	)
	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	t := &domain.Task{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, description, status, priority, project_id, assignee_id, start_date, due_date, created_by, created_at, updated_at, task_number, task_key, blocked_reason, blocked_by_task
		 FROM tasks WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID,
		&t.AssigneeID, &t.StartDate, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		&t.TaskNumber, &t.TaskKey, &t.BlockedReason, &t.BlockedByTask)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return t, err
}

type TaskFilter struct {
	ProjectID  uuid.UUID
	Status     string
	Priority   string
	AssigneeID string
	Search     string
	Page       int
	Limit      int
}

func (r *TaskRepository) List(ctx context.Context, f TaskFilter) ([]domain.Task, int, error) {
	where := "WHERE t.project_id = $1"
	args := []interface{}{f.ProjectID}
	argIdx := 2

	if f.Status != "" {
		where += fmt.Sprintf(" AND t.status = $%d", argIdx)
		args = append(args, f.Status)
		argIdx++
	}
	if f.Priority != "" {
		where += fmt.Sprintf(" AND t.priority = $%d", argIdx)
		args = append(args, f.Priority)
		argIdx++
	}
	if f.AssigneeID != "" {
		uid, err := uuid.Parse(f.AssigneeID)
		if err == nil {
			where += fmt.Sprintf(" AND t.assignee_id = $%d", argIdx)
			args = append(args, uid)
			argIdx++
		}
	}
	if f.Search != "" {
		where += fmt.Sprintf(" AND (t.title ILIKE '%%' || $%d || '%%' OR t.description ILIKE '%%' || $%d || '%%')", argIdx, argIdx)
		args = append(args, f.Search)
		argIdx++
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks t %s", where)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	query := fmt.Sprintf(
		`SELECT t.id, t.title, t.description, t.status, t.priority, t.project_id,
		        t.assignee_id, t.start_date, t.due_date, t.created_by, t.created_at, t.updated_at,
		        t.task_number, t.task_key, t.blocked_reason, t.blocked_by_task
		 FROM tasks t %s
		 ORDER BY t.created_at DESC
		 LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	args = append(args, f.Limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.ProjectID, &t.AssigneeID, &t.StartDate, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
			&t.TaskNumber, &t.TaskKey, &t.BlockedReason, &t.BlockedByTask); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}
	return tasks, total, rows.Err()
}

func (r *TaskRepository) Update(ctx context.Context, t *domain.Task) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE tasks SET title = $2, description = $3, status = $4, priority = $5,
		 assignee_id = $6, start_date = $7, due_date = $8, updated_at = $9,
		 blocked_reason = $10, blocked_by_task = $11
		 WHERE id = $1`,
		t.ID, t.Title, t.Description, t.Status, t.Priority,
		t.AssigneeID, t.StartDate, t.DueDate, t.UpdatedAt,
		t.BlockedReason, t.BlockedByTask,
	)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	return err
}

func (r *TaskRepository) CountByStatusForWorkspace(ctx context.Context, workspaceID uuid.UUID) (map[string]int, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.status, COUNT(*)
		 FROM tasks t
		 JOIN projects p ON p.id = t.project_id
		 WHERE p.workspace_id = $1
		 GROUP BY t.status`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, rows.Err()
}

func (r *TaskRepository) CountByAssigneeForWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]struct {
	UserID   uuid.UUID
	UserName string
	Count    int
}, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT u.id, u.name, COUNT(*)
		 FROM tasks t
		 JOIN projects p ON p.id = t.project_id
		 JOIN users u ON u.id = t.assignee_id
		 WHERE p.workspace_id = $1 AND t.assignee_id IS NOT NULL
		 GROUP BY u.id, u.name
		 ORDER BY COUNT(*) DESC`,
		workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		UserID   uuid.UUID
		UserName string
		Count    int
	}
	for rows.Next() {
		var item struct {
			UserID   uuid.UUID
			UserName string
			Count    int
		}
		if err := rows.Scan(&item.UserID, &item.UserName, &item.Count); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, rows.Err()
}

type MyTaskFilter struct {
	UserID     uuid.UUID
	OrgID      string
	Status     string
	Priority   string
	Search     string
	DueBefore  string
	DueAfter   string
	ProjectID  string
	Page       int
	Limit      int
}

func (r *TaskRepository) ListByAssignee(ctx context.Context, f MyTaskFilter) ([]domain.Task, int, error) {
	where := "WHERE t.assignee_id = $1"
	args := []interface{}{f.UserID}
	argIdx := 2

	if f.Status != "" {
		where += fmt.Sprintf(" AND t.status = $%d", argIdx)
		args = append(args, f.Status)
		argIdx++
	}
	if f.Priority != "" {
		where += fmt.Sprintf(" AND t.priority = $%d", argIdx)
		args = append(args, f.Priority)
		argIdx++
	}
	if f.Search != "" {
		where += fmt.Sprintf(" AND (t.title ILIKE '%%' || $%d || '%%' OR t.task_key ILIKE '%%' || $%d || '%%')", argIdx, argIdx)
		args = append(args, f.Search)
		argIdx++
	}
	if f.DueBefore != "" {
		where += fmt.Sprintf(" AND t.due_date <= $%d", argIdx)
		args = append(args, f.DueBefore)
		argIdx++
	}
	if f.DueAfter != "" {
		where += fmt.Sprintf(" AND t.due_date >= $%d", argIdx)
		args = append(args, f.DueAfter)
		argIdx++
	}
	if f.ProjectID != "" {
		uid, err := uuid.Parse(f.ProjectID)
		if err == nil {
			where += fmt.Sprintf(" AND t.project_id = $%d", argIdx)
			args = append(args, uid)
			argIdx++
		}
	}
	if f.OrgID != "" {
		orgUUID, err := uuid.Parse(f.OrgID)
		if err == nil {
			where += fmt.Sprintf(" AND t.project_id IN (SELECT p.id FROM projects p JOIN workspaces w ON p.workspace_id = w.id WHERE w.org_id = $%d)", argIdx)
			args = append(args, orgUUID)
			argIdx++
		}
	}

	var total int
	if err := r.pool.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM tasks t %s", where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	query := fmt.Sprintf(
		`SELECT t.id, t.title, t.description, t.status, t.priority, t.project_id,
		        t.assignee_id, t.start_date, t.due_date, t.created_by, t.created_at, t.updated_at,
		        t.task_number, t.task_key, t.blocked_reason, t.blocked_by_task
		 FROM tasks t %s
		 ORDER BY
		   CASE t.status WHEN 'blocked' THEN 0 WHEN 'in_progress' THEN 1 WHEN 'todo' THEN 2 ELSE 3 END,
		   COALESCE(t.due_date, '9999-12-31') ASC,
		   t.created_at DESC
		 LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, f.Limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.ProjectID, &t.AssigneeID, &t.StartDate, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
			&t.TaskNumber, &t.TaskKey, &t.BlockedReason, &t.BlockedByTask); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}
	return tasks, total, rows.Err()
}

type UserTaskStats struct {
	Total      int            `json:"total"`
	ByStatus   map[string]int `json:"by_status"`
	ByPriority map[string]int `json:"by_priority"`
	Overdue    int            `json:"overdue"`
	DueToday   int            `json:"due_today"`
	DueThisWeek int           `json:"due_this_week"`
	Completed  int            `json:"completed"`
}

func (r *TaskRepository) StatsForUser(ctx context.Context, userID uuid.UUID) (*UserTaskStats, error) {
	stats := &UserTaskStats{
		ByStatus:   make(map[string]int),
		ByPriority: make(map[string]int),
	}

	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks WHERE assignee_id = $1`, userID).Scan(&stats.Total)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT status, COUNT(*) FROM tasks WHERE assignee_id = $1 GROUP BY status`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		var c int
		if err := rows.Scan(&s, &c); err != nil {
			return nil, err
		}
		stats.ByStatus[s] = c
		if s == "done" {
			stats.Completed = c
		}
	}

	rows2, err := r.pool.Query(ctx,
		`SELECT priority, COUNT(*) FROM tasks WHERE assignee_id = $1 GROUP BY priority`, userID)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var p string
		var c int
		if err := rows2.Scan(&p, &c); err != nil {
			return nil, err
		}
		stats.ByPriority[p] = c
	}

	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks WHERE assignee_id = $1 AND due_date < CURRENT_DATE AND status != 'done'`,
		userID).Scan(&stats.Overdue); err != nil {
		return nil, err
	}

	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks WHERE assignee_id = $1 AND due_date = CURRENT_DATE AND status != 'done'`,
		userID).Scan(&stats.DueToday); err != nil {
		return nil, err
	}

	if err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tasks WHERE assignee_id = $1 AND due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days' AND status != 'done'`,
		userID).Scan(&stats.DueThisWeek); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *TaskRepository) ProjectNamesForUser(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT p.id, p.name FROM projects p
		 JOIN tasks t ON t.project_id = p.id
		 WHERE t.assignee_id = $1
		 ORDER BY p.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	m := make(map[uuid.UUID]string)
	for rows.Next() {
		var id uuid.UUID
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		m[id] = name
	}
	return m, rows.Err()
}

type OrgMemberStats struct {
	UserID    uuid.UUID      `json:"user_id"`
	UserName  string         `json:"user_name"`
	UserEmail string         `json:"user_email"`
	Role      string         `json:"role"`
	Total     int            `json:"total"`
	ByStatus  map[string]int `json:"by_status"`
	Overdue   int            `json:"overdue"`
	Completed int            `json:"completed"`
}

func (r *TaskRepository) StatsForOrgMembers(ctx context.Context, orgID uuid.UUID, viewableRoles []string) ([]OrgMemberStats, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT om.user_id, u.name, u.email, om.role,
		        COUNT(t.id) AS total,
		        COUNT(t.id) FILTER (WHERE t.status = 'todo') AS todo,
		        COUNT(t.id) FILTER (WHERE t.status = 'in_progress') AS in_progress,
		        COUNT(t.id) FILTER (WHERE t.status = 'blocked') AS blocked,
		        COUNT(t.id) FILTER (WHERE t.status = 'done') AS done,
		        COUNT(t.id) FILTER (WHERE t.due_date < CURRENT_DATE AND t.status != 'done') AS overdue
		 FROM organization_members om
		 JOIN users u ON u.id = om.user_id
		 LEFT JOIN tasks t ON t.assignee_id = om.user_id
		   AND t.project_id IN (SELECT id FROM projects WHERE org_id = $1)
		 WHERE om.org_id = $1 AND om.role = ANY($2)
		 GROUP BY om.user_id, u.name, u.email, om.role
		 ORDER BY total DESC, u.name ASC`,
		orgID, viewableRoles,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []OrgMemberStats
	for rows.Next() {
		var s OrgMemberStats
		var todo, ip, blocked, done int
		if err := rows.Scan(&s.UserID, &s.UserName, &s.UserEmail, &s.Role,
			&s.Total, &todo, &ip, &blocked, &done, &s.Overdue); err != nil {
			return nil, err
		}
		s.ByStatus = map[string]int{
			"todo": todo, "in_progress": ip, "blocked": blocked, "done": done,
		}
		s.Completed = done
		results = append(results, s)
	}
	return results, rows.Err()
}

func (r *TaskRepository) CountOverdueForWorkspace(ctx context.Context, workspaceID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM tasks t
		 JOIN projects p ON p.id = t.project_id
		 WHERE p.workspace_id = $1 AND t.due_date < CURRENT_DATE AND t.status != 'done'`,
		workspaceID,
	).Scan(&count)
	return count, err
}
