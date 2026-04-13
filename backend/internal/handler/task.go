package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/dto"
	"github.com/shivam/taskflow/backend/internal/middleware"
	"github.com/shivam/taskflow/backend/internal/repository"
	"github.com/shivam/taskflow/backend/internal/service"
	"github.com/shivam/taskflow/backend/internal/ws"
)

type TaskHandler struct {
	taskSvc    *service.TaskService
	projectSvc *service.ProjectService
	hub        *ws.Hub
}

func NewTaskHandler(taskSvc *service.TaskService, projectSvc *service.ProjectService) *TaskHandler {
	return &TaskHandler{taskSvc: taskSvc, projectSvc: projectSvc}
}

func (h *TaskHandler) SetHub(hub *ws.Hub) { h.hub = hub }

func (h *TaskHandler) broadcast(eventType string, task *domain.Task) {
	if h.hub == nil || task == nil {
		return
	}
	h.hub.BroadcastToAll(ws.Event{
		Type:      eventType,
		Payload:   task,
		ProjectID: task.ProjectID.String(),
	})
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	if _, err := h.projectSvc.VerifyWorkspaceMembership(r.Context(), projectID, userID); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	page, limit := dto.ParsePagination(r)
	filter := repository.TaskFilter{
		ProjectID:  projectID,
		Status:     r.URL.Query().Get("status"),
		Priority:   r.URL.Query().Get("priority"),
		AssigneeID: r.URL.Query().Get("assignee"),
		Search:     r.URL.Query().Get("search"),
		Page:       page,
		Limit:      limit,
	}

	tasks, total, err := h.taskSvc.List(r.Context(), filter)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if tasks == nil {
		tasks = []domain.Task{}
	}

	dto.WriteJSON(w, http.StatusOK, dto.PaginatedResponse{
		Data: tasks,
		Meta: dto.PaginationMeta{Page: page, Limit: limit, Total: total},
	})
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	if _, err := h.projectSvc.VerifyWorkspaceMembership(r.Context(), projectID, userID); err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			dto.WriteError(w, http.StatusForbidden, "forbidden")
			return
		}
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	var req dto.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	var assigneeID *uuid.UUID
	if req.AssigneeID != nil && *req.AssigneeID != "" {
		uid, err := uuid.Parse(*req.AssigneeID)
		if err != nil {
			dto.WriteValidationError(w, map[string]string{"assignee_id": "invalid uuid"})
			return
		}
		assigneeID = &uid
	}

	var startDate *time.Time
	if req.StartDate != nil {
		d, _ := time.Parse("2006-01-02", *req.StartDate)
		startDate = &d
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		d, _ := time.Parse("2006-01-02", *req.DueDate)
		dueDate = &d
	}

	task, err := h.taskSvc.Create(r.Context(), projectID, userID, req.Title, req.Description,
		domain.TaskPriority(req.Priority), assigneeID, startDate, dueDate)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	h.broadcast("task_created", task)
	dto.WriteJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	var req dto.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := req.Validate(); errs != nil {
		dto.WriteValidationError(w, errs)
		return
	}

	var status *domain.TaskStatus
	if req.Status != nil {
		s := domain.TaskStatus(*req.Status)
		status = &s
	}
	var priority *domain.TaskPriority
	if req.Priority != nil {
		p := domain.TaskPriority(*req.Priority)
		priority = &p
	}

	task, err := h.taskSvc.Update(r.Context(), taskID, userID, req.Title, req.Description, status, priority, req.AssigneeID, req.StartDate, req.DueDate, req.BlockedReason, req.BlockedByTask)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			dto.WriteError(w, http.StatusNotFound, "not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			dto.WriteError(w, http.StatusBadRequest, "invalid input")
			return
		}
		if errors.Is(err, domain.ErrConflict) {
			dto.WriteError(w, http.StatusConflict, "WIP limit exceeded")
			return
		}
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	h.broadcast("task_updated", task)
	dto.WriteJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}
	userID := middleware.GetUserID(r.Context())

	task, err := h.taskSvc.GetByID(r.Context(), taskID)
	if err != nil {
		dto.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	if err := h.taskSvc.Delete(r.Context(), taskID, userID); err != nil {
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
	if h.hub != nil {
		h.hub.BroadcastToAll(ws.Event{
			Type:      "task_deleted",
			Payload:   map[string]string{"id": taskID.String()},
			ProjectID: task.ProjectID.String(),
		})
	}
	w.WriteHeader(http.StatusNoContent)
}
