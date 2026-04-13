package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivam/taskflow/backend/internal/config"
	"github.com/shivam/taskflow/backend/internal/handler"
	"github.com/shivam/taskflow/backend/internal/middleware"
	"github.com/shivam/taskflow/backend/internal/repository"
	"github.com/shivam/taskflow/backend/internal/service"
	"github.com/shivam/taskflow/backend/internal/ws"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("connected to database")

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	workspaceRepo := repository.NewWorkspaceRepository(pool)
	projectRepo := repository.NewProjectRepository(pool)
	taskRepo := repository.NewTaskRepository(pool)
	invitationRepo := repository.NewInvitationRepository(pool)
	commentRepo := repository.NewCommentRepository(pool)
	customFieldRepo := repository.NewCustomFieldRepository(pool)
	orgRepo := repository.NewOrganizationRepository(pool)
	teamRepo := repository.NewTeamRepository(pool)
	orgInvRepo := repository.NewOrgInvitationRepository(pool)
	projectMemberRepo := repository.NewProjectMemberRepository(pool)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.BcryptCost)
	workspaceSvc := service.NewWorkspaceService(workspaceRepo, userRepo, taskRepo, projectMemberRepo)
	projectSvc := service.NewProjectService(projectRepo, workspaceRepo)
	customFieldSvc := service.NewCustomFieldService(customFieldRepo)
	taskSvc := service.NewTaskService(taskRepo, projectRepo, customFieldSvc)
	invitationSvc := service.NewInvitationService(invitationRepo, workspaceRepo, userRepo)
	commentSvc := service.NewCommentService(commentRepo)
	orgSvc := service.NewOrganizationService(orgRepo, projectMemberRepo)
	teamSvc := service.NewTeamService(teamRepo, workspaceRepo, projectMemberRepo, orgRepo)
	orgInvSvc := service.NewOrgInvitationService(orgInvRepo, orgRepo, userRepo)
	projectMemberSvc := service.NewProjectMemberService(projectMemberRepo, workspaceRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	wsHandler := handler.NewWorkspaceHandler(workspaceSvc, teamSvc, orgRepo)
	projectHandler := handler.NewProjectHandler(projectSvc)
	taskHandler := handler.NewTaskHandler(taskSvc, projectSvc)
	invitationHandler := handler.NewInvitationHandler(invitationSvc)
	commentHandler := handler.NewCommentHandler(commentSvc, projectSvc)
	customFieldHandler := handler.NewCustomFieldHandler(customFieldSvc, projectSvc, workspaceSvc)
	orgHandler := handler.NewOrganizationHandler(orgSvc, orgInvSvc)
	teamHandler := handler.NewTeamHandler(teamSvc)
	projectMemberHandler := handler.NewProjectMemberHandler(projectMemberSvc, teamSvc)
	dashboardHandler := handler.NewDashboardHandler(taskRepo)

	wsHub := ws.NewHub(logger)
	webSocketHandler := handler.NewWebSocketHandler(wsHub, cfg.JWTSecret)
	taskHandler.SetHub(wsHub)

	// Router
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))
	r.Use(chimw.RealIP)

	// Public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authSvc))

		// Personal dashboard
		r.Get("/dashboard/my-tasks", dashboardHandler.MyTasks)
		r.Get("/dashboard/my-stats", dashboardHandler.MyStats)
		r.Get("/dashboard/my-projects", dashboardHandler.MyProjects)

		// User invitations (no workspace context)
		r.Get("/invitations", invitationHandler.ListMyInvitations)
		r.Patch("/invitations/{invitationID}", invitationHandler.RespondToInvitation)

		// Org invitations for current user (top-level)
		r.Get("/org-invitations", orgHandler.ListMyOrgInvitations)
		r.Patch("/org-invitations/{invitationID}", orgHandler.RespondToOrgInvitation)

		// Organization routes
		r.Get("/organizations", orgHandler.List)
		r.Post("/organizations", orgHandler.Create)

		r.Route("/organizations/{orgID}", func(r chi.Router) {
			r.Use(middleware.OrgGuard(orgSvc))
			r.Get("/", orgHandler.Get)
			r.Patch("/", orgHandler.Update)
			r.Get("/members", orgHandler.ListMembers)
			r.Get("/prefixes", orgHandler.ListPrefixes)
			r.Get("/tasks/by-key/{taskKey}", orgHandler.GetTaskByKey)

			// Org invitations
			r.Post("/invitations", orgHandler.InviteToOrg)
			r.Get("/invitations", orgHandler.ListOrgInvitations)

			// Org member management
			r.Delete("/members/{userID}", orgHandler.RemoveMember)
			r.Post("/members/leave", orgHandler.LeaveOrg)

			// Org dashboard (manager+)
			r.Get("/dashboard/member-stats", dashboardHandler.OrgMemberStats)
			r.Get("/dashboard/member-tasks/{userID}", dashboardHandler.OrgMemberTasks)

			// Teams
			r.Get("/teams", teamHandler.List)
			r.Post("/teams", teamHandler.Create)
			r.Route("/teams/{teamID}", func(r chi.Router) {
				r.Get("/", teamHandler.Get)
				r.Patch("/", teamHandler.Update)
				r.Delete("/", teamHandler.Delete)
				r.Post("/members", teamHandler.AddMember)
				r.Delete("/members/{userID}", teamHandler.RemoveMember)
			})

			// Workspace routes (scoped to org)
			r.Get("/workspaces", wsHandler.List)
			r.Post("/workspaces", wsHandler.Create)

			r.Route("/workspaces/{workspaceID}", func(r chi.Router) {
				r.Use(middleware.WorkspaceGuard(workspaceSvc))
				r.Get("/", wsHandler.Get)
				r.Patch("/", wsHandler.Update)
				r.Delete("/", wsHandler.Delete)
				r.Get("/stats", wsHandler.Stats)
				r.Get("/members", wsHandler.ListMembers)
				r.Post("/members", invitationHandler.SendInvite)
				r.Post("/members/add", wsHandler.DirectAddMember)
				r.Post("/members/leave", wsHandler.LeaveWorkspace)
				r.Patch("/members/{userID}", wsHandler.UpdateMemberRole)
				r.Delete("/members/{userID}", wsHandler.RemoveMember)
				r.Get("/invitations", invitationHandler.ListWorkspaceInvitations)

				// Workspace teams
				r.Post("/teams", wsHandler.AddTeam)
				r.Delete("/teams/{teamID}", wsHandler.RemoveTeam)

				// Projects scoped to workspace
				r.Get("/projects", projectHandler.ListByWorkspace)
				r.Post("/projects", projectHandler.Create)

				// Project member routes
				r.Route("/projects/{projectID}", func(r chi.Router) {
					r.Get("/members", projectMemberHandler.List)
					r.Post("/members", projectMemberHandler.Add)
					r.Post("/members/leave", projectMemberHandler.Leave)
					r.Patch("/members/{userID}", projectMemberHandler.UpdateRole)
					r.Delete("/members/{userID}", projectMemberHandler.Remove)
					r.Post("/teams", projectMemberHandler.AddTeam)
				})
			})
		})

		// Direct project/task routes (verify membership inside handler)
		r.Get("/projects/{projectID}", projectHandler.Get)
		r.Patch("/projects/{projectID}", projectHandler.Update)
		r.Delete("/projects/{projectID}", projectHandler.Delete)

		// Project comments
		r.Get("/projects/{projectID}/comments", commentHandler.ListProjectComments)
		r.Post("/projects/{projectID}/comments", commentHandler.CreateProjectComment)

		// Project custom fields
		r.Get("/projects/{projectID}/custom-fields", customFieldHandler.ListDefinitions)
		r.Post("/projects/{projectID}/custom-fields", customFieldHandler.CreateDefinition)
		r.Delete("/projects/{projectID}/custom-fields/{fieldID}", customFieldHandler.DeleteDefinition)

		// Project WIP limits
		r.Get("/projects/{projectID}/wip-limits", customFieldHandler.GetWIPLimits)
		r.Put("/projects/{projectID}/wip-limits", customFieldHandler.SetWIPLimit)
		r.Delete("/projects/{projectID}/wip-limits", customFieldHandler.DeleteWIPLimit)

		// Tasks
		r.Get("/projects/{projectID}/tasks", taskHandler.List)
		r.Post("/projects/{projectID}/tasks", taskHandler.Create)

		r.Patch("/tasks/{taskID}", taskHandler.Update)
		r.Delete("/tasks/{taskID}", taskHandler.Delete)

		// Task comments
		r.Get("/tasks/{taskID}/comments", commentHandler.ListTaskComments)
		r.Post("/tasks/{taskID}/comments", commentHandler.CreateTaskComment)

		// Task custom field values
		r.Get("/tasks/{taskID}/custom-fields", customFieldHandler.GetFieldValues)
		r.Put("/tasks/{taskID}/custom-fields/{fieldID}", customFieldHandler.SetFieldValue)

		// Comment edit/delete (context-free, auth only)
		r.Patch("/comments/{commentID}", commentHandler.Update)
		r.Delete("/comments/{commentID}", commentHandler.Delete)
	})

	// WebSocket (auth via query param)
	r.Get("/ws", webSocketHandler.ServeWS)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.APIPort),
		Handler:     r,
		ReadTimeout: 15 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("server starting", "port", cfg.APIPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}
	logger.Info("server stopped")
}
