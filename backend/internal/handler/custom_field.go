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

type CustomFieldHandler struct {
	cfSvc      *service.CustomFieldService
	projectSvc *service.ProjectService
	wsSvc      *service.WorkspaceService
}

func NewCustomFieldHandler(cfSvc *service.CustomFieldService, ps *service.ProjectService, ws *service.WorkspaceService) *CustomFieldHandler {
	return &CustomFieldHandler{cfSvc: cfSvc, projectSvc: ps, wsSvc: ws}
}

// getProjectRole resolves the project and the caller's workspace role.
// Works both inside and outside the WorkspaceGuard middleware.
func (h *CustomFieldHandler) getProjectRole(r *http.Request) (uuid.UUID, domain.WorkspaceRole, error) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		return uuid.Nil, "", domain.ErrNotFound
	}
	userID := middleware.GetUserID(r.Context())

	p, err := h.projectSvc.VerifyWorkspaceMembership(r.Context(), projectID, userID)
	if err != nil {
		return uuid.Nil, "", err
	}

	role := middleware.GetWorkspaceRole(r.Context())
	if role == "" {
		member, err := h.wsSvc.GetMember(r.Context(), p.WorkspaceID, userID)
		if err != nil {
			return uuid.Nil, "", domain.ErrForbidden
		}
		role = member.Role
	}
	return p.ID, role, nil
}

func (h *CustomFieldHandler) ListDefinitions(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())
	if _, err := h.projectSvc.VerifyWorkspaceMembership(r.Context(), projectID, userID); err != nil {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	defs, err := h.cfSvc.ListDefinitions(r.Context(), projectID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if defs == nil {
		defs = []domain.CustomFieldDefinition{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"custom_fields": defs})
}

func (h *CustomFieldHandler) CreateDefinition(w http.ResponseWriter, r *http.Request) {
	projectID, role, err := h.getProjectRole(r)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.CreateCustomFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	def, err := h.cfSvc.CreateDefinition(r.Context(), projectID, req.Name, domain.CustomFieldType(req.FieldType), req.Options, req.Required, userID, role)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, def)
}

func (h *CustomFieldHandler) DeleteDefinition(w http.ResponseWriter, r *http.Request) {
	defID, err := uuid.Parse(chi.URLParam(r, "fieldID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	_, role, err := h.getProjectRole(r)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	if err := h.cfSvc.DeleteDefinition(r.Context(), defID, role); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CustomFieldHandler) SetFieldValue(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	fieldID, err := uuid.Parse(chi.URLParam(r, "fieldID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	var req dto.SetFieldValueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	val, err := h.cfSvc.SetFieldValue(r.Context(), taskID, fieldID, req.Value)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "field definition not found")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, val)
}

func (h *CustomFieldHandler) GetFieldValues(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	vals, err := h.cfSvc.GetFieldValues(r.Context(), taskID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if vals == nil {
		vals = []domain.CustomFieldValue{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"custom_fields": vals})
}

func (h *CustomFieldHandler) GetWIPLimits(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	limits, err := h.cfSvc.GetWIPLimits(r.Context(), projectID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if limits == nil {
		limits = []domain.WIPLimit{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"wip_limits": limits})
}

func (h *CustomFieldHandler) SetWIPLimit(w http.ResponseWriter, r *http.Request) {
	_, role, err := h.getProjectRole(r)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	projectID, _ := uuid.Parse(chi.URLParam(r, "projectID"))

	var req dto.SetWIPLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	limit, err := h.cfSvc.SetWIPLimit(r.Context(), projectID, domain.TaskStatus(req.Status), req.MaxTasks, role)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, limit)
}

func (h *CustomFieldHandler) DeleteWIPLimit(w http.ResponseWriter, r *http.Request) {
	_, role, err := h.getProjectRole(r)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	projectID, _ := uuid.Parse(chi.URLParam(r, "projectID"))
	status := r.URL.Query().Get("status")
	if status == "" {
		dto.WriteValidationError(w, map[string]string{"status": "query param is required"})
		return
	}

	if err := h.cfSvc.DeleteWIPLimit(r.Context(), projectID, domain.TaskStatus(status), role); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
