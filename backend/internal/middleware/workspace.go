package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/dto"
	"github.com/shivam/taskflow/backend/internal/service"
)

const WorkspaceRoleKey contextKey = "workspace_role"

func WorkspaceGuard(wsSvc *service.WorkspaceService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wsIDStr := chi.URLParam(r, "workspaceID")
			wsID, err := uuid.Parse(wsIDStr)
			if err != nil {
				dto.WriteError(w, http.StatusNotFound, "not found")
				return
			}

			userID := GetUserID(r.Context())
			member, err := wsSvc.GetMember(r.Context(), wsID, userID)
			if err != nil {
				dto.WriteError(w, http.StatusForbidden, "forbidden")
				return
			}

			ctx := context.WithValue(r.Context(), WorkspaceRoleKey, member.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetWorkspaceRole(ctx context.Context) domain.WorkspaceRole {
	role, _ := ctx.Value(WorkspaceRoleKey).(domain.WorkspaceRole)
	return role
}
