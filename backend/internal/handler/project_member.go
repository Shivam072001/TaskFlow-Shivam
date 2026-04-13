package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/dto"
	"github.com/shivam/taskflow/backend/internal/middleware"
	"github.com/shivam/taskflow/backend/internal/service"
)

type ProjectMemberHandler struct {
	pmSvc   *service.ProjectMemberService
	teamSvc *service.TeamService
}

func NewProjectMemberHandler(pmSvc *service.ProjectMemberService, teamSvc *service.TeamService) *ProjectMemberHandler {
	return &ProjectMemberHandler{pmSvc: pmSvc, teamSvc: teamSvc}
}

func (h *ProjectMemberHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	members, err := h.pmSvc.ListMembers(r.Context(), projectID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if members == nil {
		members = []domain.ProjectMember{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"members": members})
}

func (h *ProjectMemberHandler) Add(w http.ResponseWriter, r *http.Request) {
	callerRole := middleware.GetWorkspaceRole(r.Context())
	if !callerRole.CanManageMembers() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))

	var req dto.DirectAddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	member, err := h.pmSvc.AddMember(r.Context(), projectID, wsID, userID, domain.WorkspaceRole(req.Role))
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusBadRequest, "user is not a workspace member")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, member)
}

func (h *ProjectMemberHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	callerRole := middleware.GetWorkspaceRole(r.Context())
	if !callerRole.CanManageMembers() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))

	var req dto.AddTeamToEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid team_id")
		return
	}

	if err := h.teamSvc.AddTeamToProject(r.Context(), projectID, wsID, teamID, req.DefaultRole); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, map[string]string{"message": "team added to project"})
}

func (h *ProjectMemberHandler) Remove(w http.ResponseWriter, r *http.Request) {
	callerRole := middleware.GetWorkspaceRole(r.Context())
	if !callerRole.CanManageMembers() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	targetUID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	if err := h.pmSvc.RemoveMember(r.Context(), projectID, targetUID); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectMemberHandler) Leave(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	if err := h.pmSvc.RemoveMember(r.Context(), projectID, userID); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectMemberHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	callerRole := middleware.GetWorkspaceRole(r.Context())
	if !callerRole.CanManageMembers() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	targetUID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	var req dto.UpdateMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	if err := h.pmSvc.UpdateRole(r.Context(), projectID, targetUID, domain.WorkspaceRole(req.Role)); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}
