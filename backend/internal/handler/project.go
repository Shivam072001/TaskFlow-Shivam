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

type ProjectHandler struct {
	projectSvc *service.ProjectService
}

func NewProjectHandler(projectSvc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectSvc: projectSvc}
}

func (h *ProjectHandler) ListByWorkspace(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	page, limit := dto.ParsePagination(r)

	filter := repository.ProjectFilter{
		WorkspaceID: wsID,
		Search:      r.URL.Query().Get("search"),
		OwnerID:     r.URL.Query().Get("owner"),
		Page:        page,
		Limit:       limit,
	}

	projects, total, err := h.projectSvc.ListByWorkspace(r.Context(), filter)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if projects == nil {
		projects = []domain.Project{}
	}
	dto.WriteJSON(w, http.StatusOK, dto.PaginatedResponse{
		Data: projects,
		Meta: dto.PaginationMeta{Page: page, Limit: limit, Total: total},
	})
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	wsID, _ := uuid.Parse(chi.URLParam(r, "workspaceID"))
	userID := middleware.GetUserID(r.Context())

	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	project, err := h.projectSvc.Create(r.Context(), wsID, userID, req.Name, req.Prefix, req.Description)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	project, err := h.projectSvc.VerifyWorkspaceMembership(r.Context(), projectID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	project, err := h.projectSvc.Update(r.Context(), projectID, userID, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	if err := h.projectSvc.Delete(r.Context(), projectID, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
