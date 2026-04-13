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

const (
	OrgIDKey   contextKey = "org_id"
	OrgRoleKey contextKey = "org_role"
)

func OrgGuard(orgSvc *service.OrganizationService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orgIDStr := chi.URLParam(r, "orgID")
			orgID, err := uuid.Parse(orgIDStr)
			if err != nil {
				dto.WriteError(w, http.StatusNotFound, "not found")
				return
			}

			userID := GetUserID(r.Context())
			member, err := orgSvc.GetMember(r.Context(), orgID, userID)
			if err != nil {
				dto.WriteError(w, http.StatusForbidden, "forbidden")
				return
			}

			ctx := context.WithValue(r.Context(), OrgIDKey, orgID)
			ctx = context.WithValue(ctx, OrgRoleKey, member.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetOrgID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(OrgIDKey).(uuid.UUID)
	return id
}

func GetOrgRole(ctx context.Context) domain.OrgRole {
	role, _ := ctx.Value(OrgRoleKey).(domain.OrgRole)
	return role
}
