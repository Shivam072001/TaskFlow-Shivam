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

type OrganizationHandler struct {
	orgSvc    *service.OrganizationService
	orgInvSvc *service.OrgInvitationService
}

func NewOrganizationHandler(orgSvc *service.OrganizationService, orgInvSvc *service.OrgInvitationService) *OrganizationHandler {
	return &OrganizationHandler{orgSvc: orgSvc, orgInvSvc: orgInvSvc}
}

func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	orgs, err := h.orgSvc.ListByUser(r.Context(), userID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if orgs == nil {
		orgs = []domain.Organization{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"organizations": orgs})
}

func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	org, err := h.orgSvc.Create(r.Context(), userID, req.Name, req.Slug)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			dto.WriteValidationError(w, map[string]string{"slug": "already taken"})
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, org)
}

func (h *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	org, err := h.orgSvc.GetByID(r.Context(), orgID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, org)
}

func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	role := middleware.GetOrgRole(r.Context())

	var req dto.UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	org, err := h.orgSvc.Update(r.Context(), orgID, role, req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, org)
}

func (h *OrganizationHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	members, err := h.orgSvc.ListMembers(r.Context(), orgID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if members == nil {
		members = []domain.OrgMember{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"members": members})
}

func (h *OrganizationHandler) ListPrefixes(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	prefixes, err := h.orgSvc.ListPrefixes(r.Context(), orgID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if prefixes == nil {
		prefixes = []string{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"prefixes": prefixes})
}

func (h *OrganizationHandler) GetTaskByKey(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	taskKey := chi.URLParam(r, "taskKey")

	task, workspaceID, err := h.orgSvc.GetTaskByKey(r.Context(), orgID, taskKey)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "task not found")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	type taskWithContext struct {
		*domain.Task
		WorkspaceID uuid.UUID `json:"workspace_id"`
	}

	dto.WriteJSON(w, http.StatusOK, taskWithContext{Task: task, WorkspaceID: workspaceID})
}

func (h *OrganizationHandler) InviteToOrg(w http.ResponseWriter, r *http.Request) {
	orgRole := middleware.GetOrgRole(r.Context())
	if !orgRole.CanInviteToOrg() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	var req dto.SendOrgInviteRequest
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

	inv, err := h.orgInvSvc.SendInvite(r.Context(), orgID, userID, req.Email, domain.OrgRole(req.Role))
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			dto.WriteValidationError(w, map[string]string{"email": "user is already a member"})
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, inv)
}

func (h *OrganizationHandler) ListOrgInvitations(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	invitations, err := h.orgInvSvc.ListByOrg(r.Context(), orgID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if invitations == nil {
		invitations = []domain.OrgInvitation{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"invitations": invitations})
}

func (h *OrganizationHandler) ListMyOrgInvitations(w http.ResponseWriter, r *http.Request) {
	email := middleware.GetEmail(r.Context())
	invitations, err := h.orgInvSvc.ListByEmail(r.Context(), email)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if invitations == nil {
		invitations = []domain.OrgInvitation{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"invitations": invitations})
}

func (h *OrganizationHandler) RespondToOrgInvitation(w http.ResponseWriter, r *http.Request) {
	invID, err := uuid.Parse(chi.URLParam(r, "invitationID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	var req dto.RespondOrgInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	userID := middleware.GetUserID(r.Context())
	accept := req.Action == "accept"

	if err := h.orgInvSvc.Respond(r.Context(), invID, userID, accept); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			dto.WriteError(w, http.StatusConflict, "invitation already responded")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	verb := "accepted"
	if req.Action == "decline" {
		verb = "declined"
	}
	dto.WriteJSON(w, http.StatusOK, map[string]string{"message": "invitation " + verb})
}

func (h *OrganizationHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	orgRole := middleware.GetOrgRole(r.Context())
	if !orgRole.CanManageOrg() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	targetUID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	orgID := middleware.GetOrgID(r.Context())
	if err := h.orgSvc.RemoveMemberCascade(r.Context(), orgID, targetUID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrganizationHandler) LeaveOrg(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	userID := middleware.GetUserID(r.Context())

	if err := h.orgSvc.RemoveMemberCascade(r.Context(), orgID, userID); err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
