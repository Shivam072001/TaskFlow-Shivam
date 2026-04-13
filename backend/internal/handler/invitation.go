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

type InvitationHandler struct {
	invitationSvc *service.InvitationService
}

func NewInvitationHandler(invSvc *service.InvitationService) *InvitationHandler {
	return &InvitationHandler{invitationSvc: invSvc}
}

func (h *InvitationHandler) SendInvite(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	callerID := middleware.GetUserID(r.Context())
	callerRole := middleware.GetWorkspaceRole(r.Context())

	var req dto.InviteMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	inv, err := h.invitationSvc.SendInvite(r.Context(), wsID, callerID, callerRole, req.Email, domain.WorkspaceRole(req.Role))
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			dto.WriteValidationError(w, map[string]string{"email": "already invited or already a member"})
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			dto.WriteError(w, http.StatusBadRequest, "invalid role")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, inv)
}

func (h *InvitationHandler) ListMyInvitations(w http.ResponseWriter, r *http.Request) {
	email := middleware.GetEmail(r.Context())

	invitations, err := h.invitationSvc.ListPending(r.Context(), email)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if invitations == nil {
		invitations = []domain.WorkspaceInvitation{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"invitations": invitations})
}

func (h *InvitationHandler) ListWorkspaceInvitations(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))

	invitations, err := h.invitationSvc.ListByWorkspace(r.Context(), wsID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if invitations == nil {
		invitations = []domain.WorkspaceInvitation{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"invitations": invitations})
}

func (h *InvitationHandler) RespondToInvitation(w http.ResponseWriter, r *http.Request) {
	invID, err := uuid.Parse(chi.URLParam(r, "invitationID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())
	email := middleware.GetEmail(r.Context())

	var req dto.InvitationResponseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	accept := req.Action == "accept"
	if err := h.invitationSvc.Respond(r.Context(), invID, userID, email, accept); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			dto.WriteError(w, http.StatusConflict, "invitation already responded to")
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
