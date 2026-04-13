package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/domain"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) Create(ctx context.Context, t *domain.Team) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO teams (id, org_id, name, created_by, created_at) VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.OrgID, t.Name, t.CreatedBy, t.CreatedAt,
	)
	return err
}

func (r *TeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	t := &domain.Team{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, org_id, name, created_by, created_at FROM teams WHERE id = $1`, id,
	).Scan(&t.ID, &t.OrgID, &t.Name, &t.CreatedBy, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return t, err
}

func (r *TeamRepository) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.Team, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, org_id, name, created_by, created_at FROM teams WHERE org_id = $1 ORDER BY name`, orgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.CreatedBy, &t.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *TeamRepository) Update(ctx context.Context, id uuid.UUID, name string) error {
	_, err := r.pool.Exec(ctx, `UPDATE teams SET name = $2 WHERE id = $1`, id, name)
	return err
}

func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM teams WHERE id = $1`, id)
	return err
}

func (r *TeamRepository) AddMember(ctx context.Context, m *domain.TeamMember) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO team_members (id, team_id, user_id, added_at) VALUES ($1, $2, $3, $4)
		 ON CONFLICT (team_id, user_id) DO NOTHING`,
		m.ID, m.TeamID, m.UserID, m.AddedAt,
	)
	return err
}

func (r *TeamRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`, teamID, userID,
	)
	return err
}

func (r *TeamRepository) ListMembers(ctx context.Context, teamID uuid.UUID) ([]domain.TeamMember, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT tm.id, tm.team_id, tm.user_id, tm.added_at, u.name, u.email
		 FROM team_members tm
		 JOIN users u ON u.id = tm.user_id
		 WHERE tm.team_id = $1
		 ORDER BY tm.added_at`, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.TeamMember
	for rows.Next() {
		var m domain.TeamMember
		if err := rows.Scan(&m.ID, &m.TeamID, &m.UserID, &m.AddedAt, &m.UserName, &m.UserEmail); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

// AddTeamToWorkspace records the team assignment and returns the workspace_teams ID.
func (r *TeamRepository) AddTeamToWorkspace(ctx context.Context, wt *domain.WorkspaceTeam) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO workspace_teams (id, workspace_id, team_id, default_role, added_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (workspace_id, team_id) DO UPDATE SET default_role = EXCLUDED.default_role`,
		wt.ID, wt.WorkspaceID, wt.TeamID, wt.DefaultRole, wt.AddedAt,
	)
	return err
}

func (r *TeamRepository) RemoveTeamFromWorkspace(ctx context.Context, workspaceID, teamID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM workspace_teams WHERE workspace_id = $1 AND team_id = $2`, workspaceID, teamID,
	)
	return err
}

func (r *TeamRepository) AddTeamToProject(ctx context.Context, pt *domain.ProjectTeam) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO project_teams (id, project_id, team_id, default_role, added_at)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (project_id, team_id) DO UPDATE SET default_role = EXCLUDED.default_role`,
		pt.ID, pt.ProjectID, pt.TeamID, pt.DefaultRole, pt.AddedAt,
	)
	return err
}

func (r *TeamRepository) RemoveTeamFromProject(ctx context.Context, projectID, teamID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM project_teams WHERE project_id = $1 AND team_id = $2`, projectID, teamID,
	)
	return err
}

// ListWorkspaceTeamAssignments returns all workspace_teams rows for a given team.
func (r *TeamRepository) ListWorkspaceTeamAssignments(ctx context.Context, teamID uuid.UUID) ([]domain.WorkspaceTeam, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workspace_id, team_id, default_role, added_at FROM workspace_teams WHERE team_id = $1`, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wts []domain.WorkspaceTeam
	for rows.Next() {
		var wt domain.WorkspaceTeam
		if err := rows.Scan(&wt.ID, &wt.WorkspaceID, &wt.TeamID, &wt.DefaultRole, &wt.AddedAt); err != nil {
			return nil, err
		}
		wts = append(wts, wt)
	}
	return wts, rows.Err()
}

// ListProjectTeamAssignments returns all project_teams rows for a given team.
func (r *TeamRepository) ListProjectTeamAssignments(ctx context.Context, teamID uuid.UUID) ([]domain.ProjectTeam, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, project_id, team_id, default_role, added_at FROM project_teams WHERE team_id = $1`, teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pts []domain.ProjectTeam
	for rows.Next() {
		var pt domain.ProjectTeam
		if err := rows.Scan(&pt.ID, &pt.ProjectID, &pt.TeamID, &pt.DefaultRole, &pt.AddedAt); err != nil {
			return nil, err
		}
		pts = append(pts, pt)
	}
	return pts, rows.Err()
}

// ListTeamsByWorkspace returns teams assigned to a workspace.
func (r *TeamRepository) ListTeamsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]domain.Team, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.org_id, t.name, t.created_by, t.created_at
		 FROM teams t
		 JOIN workspace_teams wt ON wt.team_id = t.id
		 WHERE wt.workspace_id = $1
		 ORDER BY t.name`, workspaceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.CreatedBy, &t.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

// ListTeamsByProject returns teams assigned to a project.
func (r *TeamRepository) ListTeamsByProject(ctx context.Context, projectID uuid.UUID) ([]domain.Team, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.org_id, t.name, t.created_by, t.created_at
		 FROM teams t
		 JOIN project_teams pt ON pt.team_id = t.id
		 WHERE pt.project_id = $1
		 ORDER BY t.name`, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var t domain.Team
		if err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.CreatedBy, &t.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}
