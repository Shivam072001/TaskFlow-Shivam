package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), 12)
	if err != nil {
		log.Fatalf("bcrypt: %v", err)
	}
	pw := string(hash)

	now := time.Now().UTC()
	uid := make([]uuid.UUID, 10)
	for i := range uid {
		uid[i] = uuid.New()
	}

	users := []struct {
		name, email string
	}{
		{"Alice Owner", "test@example.com"},
		{"Bob Admin", "bob@example.com"},
		{"Charlie Manager", "charlie@example.com"},
		{"Diana Lead", "diana@example.com"},
		{"Eve Developer", "eve@example.com"},
		{"Frank Designer", "frank@example.com"},
		{"Grace Tester", "grace@example.com"},
		{"Hank Analyst", "hank@example.com"},
		{"Ivy Intern", "ivy@example.com"},
		{"Jack Viewer", "jack@example.com"},
	}

	log.Println("Seeding users...")
	for i, u := range users {
		_, err := pool.Exec(ctx,
			`INSERT INTO users (id, name, email, password, created_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (email) DO NOTHING`,
			uid[i], u.name, u.email, pw, now)
		if err != nil {
			log.Fatalf("user %d: %v", i, err)
		}
	}

	// 3 organisations – Alice owns all
	orgIDs := [3]uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	orgs := []struct {
		name, slug string
	}{
		{"Acme Corp", "acme-corp"},
		{"Globex Inc", "globex-inc"},
		{"Initech", "initech"},
	}

	log.Println("Seeding organisations...")
	for i, o := range orgs {
		_, err := pool.Exec(ctx,
			`INSERT INTO organizations (id, name, slug, created_by, created_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (slug) DO NOTHING`,
			orgIDs[i], o.name, o.slug, uid[0], now)
		if err != nil {
			log.Fatalf("org %d: %v", i, err)
		}
	}

	// Org members – Alice=owner everywhere, spread others across orgs
	orgMembers := []struct {
		org  int
		user int
		role string
	}{
		{0, 0, "owner"}, {0, 1, "admin"}, {0, 2, "manager"}, {0, 3, "member"},
		{0, 4, "member"}, {0, 5, "member"}, {0, 6, "member"},
		{1, 0, "owner"}, {1, 7, "admin"}, {1, 8, "member"},
		{2, 0, "owner"}, {2, 9, "admin"}, {2, 3, "member"}, {2, 4, "member"},
	}

	log.Println("Seeding org members...")
	for _, m := range orgMembers {
		_, _ = pool.Exec(ctx,
			`INSERT INTO organization_members (id, org_id, user_id, role, joined_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (org_id, user_id) DO NOTHING`,
			uuid.New(), orgIDs[m.org], uid[m.user], m.role, now)
	}

	// 3 workspaces – spread across orgs
	wsIDs := [3]uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	workspaces := []struct {
		name, desc string
		org        int
	}{
		{"Engineering", "Core product development", 0},
		{"Marketing", "Campaigns and content", 1},
		{"Operations", "Internal tooling", 2},
	}

	log.Println("Seeding workspaces...")
	for i, w := range workspaces {
		_, err := pool.Exec(ctx,
			`INSERT INTO workspaces (id, org_id, name, description, created_by, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7) ON CONFLICT DO NOTHING`,
			wsIDs[i], orgIDs[w.org], w.name, w.desc, uid[0], now, now)
		if err != nil {
			log.Fatalf("ws %d: %v", i, err)
		}
	}

	// Workspace members
	wsMembers := []struct {
		ws   int
		user int
		role string
	}{
		{0, 0, "owner"}, {0, 1, "admin"}, {0, 2, "manager"}, {0, 3, "lead"},
		{0, 4, "member"}, {0, 5, "member"}, {0, 6, "member"},
		{1, 0, "owner"}, {1, 7, "admin"}, {1, 8, "member"},
		{2, 0, "owner"}, {2, 9, "admin"}, {2, 3, "member"}, {2, 4, "member"},
	}

	log.Println("Seeding workspace members...")
	for _, m := range wsMembers {
		_, _ = pool.Exec(ctx,
			`INSERT INTO workspace_members (id, workspace_id, user_id, role, joined_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (workspace_id, user_id) DO NOTHING`,
			uuid.New(), wsIDs[m.ws], uid[m.user], m.role, now)
	}

	// 2 teams in org 0
	teamIDs := [2]uuid.UUID{uuid.New(), uuid.New()}
	teams := []struct {
		name string
	}{
		{"Backend Squad"},
		{"Frontend Squad"},
	}

	log.Println("Seeding teams...")
	for i, t := range teams {
		_, err := pool.Exec(ctx,
			`INSERT INTO teams (id, org_id, name, created_by, created_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (org_id, name) DO NOTHING`,
			teamIDs[i], orgIDs[0], t.name, uid[0], now)
		if err != nil {
			log.Fatalf("team %d: %v", i, err)
		}
	}

	// Team members
	teamMembers := []struct {
		team int
		user int
	}{
		{0, 2}, {0, 3}, {0, 4},
		{1, 5}, {1, 6},
	}
	for _, m := range teamMembers {
		_, _ = pool.Exec(ctx,
			`INSERT INTO team_members (id, team_id, user_id, added_at) VALUES ($1,$2,$3,$4) ON CONFLICT (team_id, user_id) DO NOTHING`,
			uuid.New(), teamIDs[m.team], uid[m.user], now)
	}

	// 10 projects – spread across workspaces, mostly ws 0 for pagination testing
	projIDs := make([]uuid.UUID, 10)
	for i := range projIDs {
		projIDs[i] = uuid.New()
	}

	projects := []struct {
		name, prefix, desc string
		ws                 int
	}{
		{"API Gateway", "APIGW", "REST API gateway service", 0},
		{"User Service", "USER", "Authentication and profiles", 0},
		{"Payment Engine", "PAY", "Billing and payments", 0},
		{"Notification Hub", "NOTIF", "Email, SMS, push", 0},
		{"Admin Dashboard", "ADMIN", "Internal admin panel", 0},
		{"Mobile App", "MOB", "iOS and Android app", 0},
		{"Data Pipeline", "DATA", "ETL and analytics", 0},
		{"Marketing Site", "MKTG", "Public website", 1},
		{"CRM Integration", "CRM", "Salesforce sync", 1},
		{"Internal Wiki", "WIKI", "Knowledge base", 2},
	}

	log.Println("Seeding projects...")
	for i, p := range projects {
		_, err := pool.Exec(ctx,
			`INSERT INTO projects (id, org_id, name, prefix, description, workspace_id, owner_id, created_at)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT DO NOTHING`,
			projIDs[i], orgIDs[projects[i].ws], p.name, p.prefix, p.desc, wsIDs[p.ws], uid[0], now)
		if err != nil {
			log.Fatalf("project %d: %v", i, err)
		}
	}

	// Tasks – 4-6 per status per project (todo, in_progress, blocked, done)
	statuses := []string{"todo", "in_progress", "blocked", "done"}
	priorities := []string{"low", "medium", "high"}
	taskVerbs := []string{"Implement", "Fix", "Refactor", "Test", "Review", "Deploy"}
	taskNouns := []string{"auth flow", "dashboard", "API endpoint", "search", "cache layer", "logging"}

	log.Println("Seeding tasks...")
	taskCounter := make(map[uuid.UUID]int)
	for pi, projID := range projIDs {
		for _, status := range statuses {
			count := 4 + (pi+len(status))%3 // 4, 5, or 6 per status
			for t := 0; t < count; t++ {
				taskCounter[projID]++
				num := taskCounter[projID]
				verb := taskVerbs[(pi+t)%len(taskVerbs)]
				noun := taskNouns[(pi+t+num)%len(taskNouns)]
				title := fmt.Sprintf("%s %s", verb, noun)
				prio := priorities[(pi+t)%len(priorities)]
				assignee := uid[(pi+t)%len(uid)]
				prefix := projects[pi].prefix
				key := fmt.Sprintf("%s-%03d", prefix, num)

				blockedReason := ""
				blockedByTask := ""
				if status == "blocked" {
					blockedReason = "Waiting on dependency"
					if num > 1 {
						blockedByTask = fmt.Sprintf("%s-%03d", prefix, num-1)
					}
				}

				_, err := pool.Exec(ctx,
					`INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, start_date, due_date, created_by, created_at, updated_at, task_number, task_key, blocked_reason, blocked_by_task)
					 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) ON CONFLICT DO NOTHING`,
					uuid.New(), title, fmt.Sprintf("Description for: %s", title), status, prio,
					projID, assignee,
					now.AddDate(0, 0, -5), now.AddDate(0, 0, 10+t),
					uid[0], now, now,
					num, key, blockedReason, blockedByTask,
				)
				if err != nil {
					log.Fatalf("task p%d s=%s t=%d: %v", pi, status, t, err)
				}
			}
		}
	}

	// Project members for each project (subset of workspace members)
	log.Println("Seeding project members...")
	projMemberSets := [][]struct {
		user int
		role string
	}{
		// ws0 projects (0-6): pick from ws0 members
		{{0, "owner"}, {1, "admin"}, {2, "manager"}, {3, "lead"}, {4, "member"}},
		{{0, "owner"}, {2, "admin"}, {4, "member"}, {5, "member"}},
		{{0, "owner"}, {1, "admin"}, {6, "member"}},
		{{0, "owner"}, {3, "lead"}, {5, "member"}, {6, "member"}},
		{{0, "owner"}, {1, "admin"}, {2, "manager"}},
		{{0, "owner"}, {4, "member"}, {5, "member"}, {6, "member"}},
		{{0, "owner"}, {2, "manager"}, {3, "lead"}},
		// ws1 projects (7-8)
		{{0, "owner"}, {7, "admin"}, {8, "member"}},
		{{0, "owner"}, {7, "admin"}},
		// ws2 project (9)
		{{0, "owner"}, {9, "admin"}, {3, "member"}, {4, "member"}},
	}

	for pi, members := range projMemberSets {
		for _, m := range members {
			_, _ = pool.Exec(ctx,
				`INSERT INTO project_members (id, project_id, user_id, role, joined_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (project_id, user_id) DO NOTHING`,
				uuid.New(), projIDs[pi], uid[m.user], m.role, now)
		}
	}

	// Assign teams to workspace 0
	log.Println("Seeding workspace/project team assignments...")
	for _, tid := range teamIDs {
		_, _ = pool.Exec(ctx,
			`INSERT INTO workspace_teams (id, workspace_id, team_id, default_role, added_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (workspace_id, team_id) DO NOTHING`,
			uuid.New(), wsIDs[0], tid, "member", now)
	}

	// Assign Backend Squad to first 3 projects, Frontend Squad to projects 4-6
	for i := 0; i < 3; i++ {
		_, _ = pool.Exec(ctx,
			`INSERT INTO project_teams (id, project_id, team_id, default_role, added_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (project_id, team_id) DO NOTHING`,
			uuid.New(), projIDs[i], teamIDs[0], "member", now)
	}
	for i := 3; i < 6; i++ {
		_, _ = pool.Exec(ctx,
			`INSERT INTO project_teams (id, project_id, team_id, default_role, added_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (project_id, team_id) DO NOTHING`,
			uuid.New(), projIDs[i], teamIDs[1], "member", now)
	}

	log.Println("Seed complete!")
	log.Println("")
	log.Println("Test credentials:")
	log.Println("  Email:    test@example.com")
	log.Println("  Password: password123")
	log.Println("")
	log.Printf("Created: 10 users, 3 orgs, 3 workspaces, 2 teams, 10 projects, ~%d tasks\n",
		sumValues(taskCounter))
}

func sumValues(m map[uuid.UUID]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}
