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

type TeamHandler struct {
	teamSvc *service.TeamService
}

func NewTeamHandler(ts *service.TeamService) *TeamHandler {
	return &TeamHandler{teamSvc: ts}
}

func (h *TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	teams, err := h.teamSvc.ListByOrg(r.Context(), orgID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if teams == nil {
		teams = []domain.Team{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"teams": teams})
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgRole := middleware.GetOrgRole(r.Context())
	if !orgRole.CanManageTeams() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req dto.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	orgID := middleware.GetOrgID(r.Context())
	userID := middleware.GetUserID(r.Context())

	team, err := h.teamSvc.Create(r.Context(), orgID, userID, req.Name)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, team)
}

func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	teamID, err := uuid.Parse(chi.URLParam(r, "teamID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	team, err := h.teamSvc.GetByID(r.Context(), teamID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	members, err := h.teamSvc.ListMembers(r.Context(), teamID)
	if err != nil {
		members = []domain.TeamMember{}
	}
	if members == nil {
		members = []domain.TeamMember{}
	}

	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"team":    team,
		"members": members,
	})
}

func (h *TeamHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgRole := middleware.GetOrgRole(r.Context())
	if !orgRole.CanManageTeams() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	teamID, err := uuid.Parse(chi.URLParam(r, "teamID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	var req dto.UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	if err := h.teamSvc.Update(r.Context(), teamID, req.Name); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *TeamHandler) Delete(w http.ResponseWriter, r *http.Request) {
	orgRole := middleware.GetOrgRole(r.Context())
	if !orgRole.CanManageTeams() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	teamID, err := uuid.Parse(chi.URLParam(r, "teamID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	if err := h.teamSvc.Delete(r.Context(), teamID); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	orgRole := middleware.GetOrgRole(r.Context())
	if !orgRole.CanManageTeams() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	teamID, err := uuid.Parse(chi.URLParam(r, "teamID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	var req dto.AddTeamMemberRequest
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

	member, err := h.teamSvc.AddMember(r.Context(), teamID, userID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, member)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	orgRole := middleware.GetOrgRole(r.Context())
	if !orgRole.CanManageTeams() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	teamID, err := uuid.Parse(chi.URLParam(r, "teamID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	if err := h.teamSvc.RemoveMember(r.Context(), teamID, userID); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
