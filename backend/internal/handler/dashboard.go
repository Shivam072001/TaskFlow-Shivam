package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/dto"
	"github.com/shivam/taskflow/backend/internal/middleware"
	"github.com/shivam/taskflow/backend/internal/repository"
)

var rolePriority = map[domain.OrgRole]int{
	domain.OrgRoleOwner:   4,
	domain.OrgRoleAdmin:   3,
	domain.OrgRoleManager: 2,
	domain.OrgRoleMember:  1,
}

func viewableRolesBelow(callerRole domain.OrgRole) []string {
	callerLevel := rolePriority[callerRole]
	var roles []string
	for role, level := range rolePriority {
		if level < callerLevel {
			roles = append(roles, string(role))
		}
	}
	return roles
}

type DashboardHandler struct {
	taskRepo *repository.TaskRepository
}

func NewDashboardHandler(taskRepo *repository.TaskRepository) *DashboardHandler {
	return &DashboardHandler{taskRepo: taskRepo}
}

func (h *DashboardHandler) MyTasks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	page, limit := dto.ParsePagination(r)
	q := r.URL.Query()

	filter := repository.MyTaskFilter{
		UserID:    userID,
		Status:    q.Get("status"),
		Priority:  q.Get("priority"),
		Search:    q.Get("search"),
		DueBefore: q.Get("due_before"),
		DueAfter:  q.Get("due_after"),
		ProjectID: q.Get("project_id"),
		Page:      page,
		Limit:     limit,
	}

	tasks, total, err := h.taskRepo.ListByAssignee(r.Context(), filter)
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

func (h *DashboardHandler) MyStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	stats, err := h.taskRepo.StatsForUser(r.Context(), userID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	dto.WriteJSON(w, http.StatusOK, stats)
}

func (h *DashboardHandler) MyProjects(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	projects, err := h.taskRepo.ProjectNamesForUser(r.Context(), userID)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	type item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var list []item
	for id, name := range projects {
		list = append(list, item{ID: id.String(), Name: name})
	}
	if list == nil {
		list = []item{}
	}
	dto.WriteJSON(w, http.StatusOK, list)
}

func (h *DashboardHandler) OrgMemberStats(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	callerRole := middleware.GetOrgRole(r.Context())

	if !callerRole.CanManageTeams() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	roles := viewableRolesBelow(callerRole)
	if len(roles) == 0 {
		dto.WriteJSON(w, http.StatusOK, []interface{}{})
		return
	}

	stats, err := h.taskRepo.StatsForOrgMembers(r.Context(), orgID, roles)
	if err != nil {
		dto.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if stats == nil {
		stats = []repository.OrgMemberStats{}
	}
	dto.WriteJSON(w, http.StatusOK, stats)
}

func (h *DashboardHandler) OrgMemberTasks(w http.ResponseWriter, r *http.Request) {
	orgID := middleware.GetOrgID(r.Context())
	callerRole := middleware.GetOrgRole(r.Context())

	if !callerRole.CanManageTeams() {
		dto.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	targetUserID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		dto.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	page, limit := dto.ParsePagination(r)
	q := r.URL.Query()
	filter := repository.MyTaskFilter{
		UserID:    targetUserID,
		OrgID:     orgID.String(),
		Status:    q.Get("status"),
		Priority:  q.Get("priority"),
		Search:    q.Get("search"),
		DueBefore: q.Get("due_before"),
		DueAfter:  q.Get("due_after"),
		ProjectID: q.Get("project_id"),
		Page:      page,
		Limit:     limit,
	}

	tasks, total, err := h.taskRepo.ListByAssignee(r.Context(), filter)
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
