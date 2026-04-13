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
	"github.com/shivam/taskflow/backend/internal/repository"
	"github.com/shivam/taskflow/backend/internal/service"
)

type WorkspaceHandler struct {
	wsSvc   *service.WorkspaceService
	teamSvc *service.TeamService
	orgRepo *repository.OrganizationRepository
}

func NewWorkspaceHandler(wsSvc *service.WorkspaceService, teamSvc *service.TeamService, orgRepo *repository.OrganizationRepository) *WorkspaceHandler {
	return &WorkspaceHandler{wsSvc: wsSvc, teamSvc: teamSvc, orgRepo: orgRepo}
}

func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	orgID := middleware.GetOrgID(r.Context())
	page, limit := dto.ParsePagination(r)

	workspaces, total, err := h.wsSvc.ListByUser(r.Context(), userID, orgID, page, limit)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if workspaces == nil {
		workspaces = []domain.Workspace{}
	}
	dto.WriteJSON(w, http.StatusOK, dto.PaginatedResponse{
		Data: workspaces,
		Meta: dto.PaginationMeta{Page: page, Limit: limit, Total: total},
	})
}

func (h *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	userID := middleware.GetUserID(r.Context())
	orgID := middleware.GetOrgID(r.Context())
	ws, err := h.wsSvc.Create(r.Context(), orgID, userID, req.Name, req.Description)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, ws)
}

func (h *WorkspaceHandler) Get(w http.ResponseWriter, r *http.Request) {
	wsID, err := uuid.Parse(chi.URLParam(r, "workspaceID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	ws, err := h.wsSvc.GetByID(r.Context(), wsID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	role := middleware.GetWorkspaceRole(r.Context())

	var req dto.UpdateWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	ws, err := h.wsSvc.Update(r.Context(), wsID, role, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, ws)
}

func (h *WorkspaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	role := middleware.GetWorkspaceRole(r.Context())

	if err := h.wsSvc.Delete(r.Context(), wsID, role); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))

	members, err := h.wsSvc.ListMembers(r.Context(), wsID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if members == nil {
		members = []domain.WorkspaceMember{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"members": members})
}

func (h *WorkspaceHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	targetUID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	callerRole := middleware.GetWorkspaceRole(r.Context())

	var req dto.UpdateMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	if err := h.wsSvc.UpdateMemberRole(r.Context(), wsID, targetUID, callerRole, domain.WorkspaceRole(req.Role)); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}

func (h *WorkspaceHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	targetUID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	callerID := middleware.GetUserID(r.Context())
	callerRole := middleware.GetWorkspaceRole(r.Context())

	if err := h.wsSvc.RemoveMember(r.Context(), wsID, targetUID, callerID, callerRole); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) Stats(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))

	stats, err := h.wsSvc.GetStats(r.Context(), wsID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, stats)
}

func (h *WorkspaceHandler) DirectAddMember(w http.ResponseWriter, r *http.Request) {
	callerRole := middleware.GetWorkspaceRole(r.Context())
	if !callerRole.CanManageMembers() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	orgID := middleware.GetOrgID(r.Context())

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

	// Validate user is an org member
	_, err = h.orgRepo.GetMember(r.Context(), orgID, userID)
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "user is not an org member")
		return
	}

	if err := h.wsSvc.DirectAddMember(r.Context(), wsID, userID, domain.WorkspaceRole(req.Role)); err != nil {
		if errors.Is(err, domain.ErrConflict) {
			dto.WriteError(w, http.StatusConflict, "user is already a workspace member")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, map[string]string{"message": "member added"})
}

func (h *WorkspaceHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	callerRole := middleware.GetWorkspaceRole(r.Context())
	if !callerRole.CanManageMembers() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
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

	if err := h.teamSvc.AddTeamToWorkspace(r.Context(), wsID, teamID, req.DefaultRole); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, map[string]string{"message": "team added to workspace"})
}

func (h *WorkspaceHandler) RemoveTeam(w http.ResponseWriter, r *http.Request) {
	callerRole := middleware.GetWorkspaceRole(r.Context())
	if !callerRole.CanManageMembers() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	teamID, err := uuid.Parse(chi.URLParam(r, "teamID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	if err := h.teamSvc.RemoveTeamFromWorkspace(r.Context(), wsID, teamID); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkspaceHandler) LeaveWorkspace(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	userID := middleware.GetUserID(r.Context())
	callerRole := middleware.GetWorkspaceRole(r.Context())

	if err := h.wsSvc.RemoveMember(r.Context(), wsID, userID, userID, callerRole); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
