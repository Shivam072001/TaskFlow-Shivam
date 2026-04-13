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

type CommentHandler struct {
	commentSvc *service.CommentService
	projectSvc *service.ProjectService
}

func NewCommentHandler(cs *service.CommentService, ps *service.ProjectService) *CommentHandler {
	return &CommentHandler{commentSvc: cs, projectSvc: ps}
}

func (h *CommentHandler) ListProjectComments(w http.ResponseWriter, r *http.Request) {
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

	comments, err := h.commentSvc.List(r.Context(), domain.EntityProject, projectID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if comments == nil {
		comments = []domain.Comment{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"comments": comments})
}

func (h *CommentHandler) CreateProjectComment(w http.ResponseWriter, r *http.Request) {
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

	h.createComment(w, r, domain.EntityProject, projectID, userID)
}

func (h *CommentHandler) ListTaskComments(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	comments, err := h.commentSvc.List(r.Context(), domain.EntityTask, taskID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if comments == nil {
		comments = []domain.Comment{}
	}
	dto.WriteJSON(w, http.StatusOK, map[string]interface{}{"comments": comments})
}

func (h *CommentHandler) CreateTaskComment(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	h.createComment(w, r, domain.EntityTask, taskID, userID)
}

func (h *CommentHandler) createComment(w http.ResponseWriter, r *http.Request, entityType domain.CommentEntityType, entityID, userID uuid.UUID) {
	var req struct {
		ParentID *string `json:"parent_id"`
		Content  string  `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Content == "" {
		dto.WriteValidationError(w, map[string]string{"content": "is required"})
		return
	}

	var parentID *uuid.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			dto.WriteValidationError(w, map[string]string{"parent_id": "invalid uuid"})
			return
		}
		parentID = &pid
	}

	comment, err := h.commentSvc.Create(r.Context(), entityType, entityID, userID, parentID, req.Content)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusCreated, comment)
}

func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	comment, err := h.commentSvc.Update(r.Context(), commentID, userID, req.Content)
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
	dto.WriteJSON(w, http.StatusOK, comment)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	commentID, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())
	role := middleware.GetWorkspaceRole(r.Context())

	if err := h.commentSvc.Delete(r.Context(), commentID, userID, role); err != nil {
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
