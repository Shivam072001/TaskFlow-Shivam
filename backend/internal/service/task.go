package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type TaskService struct {
	taskRepo    *repository.TaskRepository
	projectRepo *repository.ProjectRepository
	cfSvc       *CustomFieldService
}

func NewTaskService(tr *repository.TaskRepository, pr *repository.ProjectRepository, cfSvc *CustomFieldService) *TaskService {
	return &TaskService{taskRepo: tr, projectRepo: pr, cfSvc: cfSvc}
}

func (s *TaskService) Create(ctx context.Context, projectID, creatorID uuid.UUID, title, description string, priority domain.TaskPriority, assigneeID *uuid.UUID, startDate, dueDate *time.Time) (*domain.Task, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("fetching project for task numbering: %w", err)
	}

	num, err := s.taskRepo.NextTaskNumber(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("generating task number: %w", err)
	}

	now := time.Now().UTC()
	t := &domain.Task{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      domain.StatusTodo,
		Priority:    priority,
		ProjectID:   projectID,
		AssigneeID:  assigneeID,
		StartDate:   startDate,
		DueDate:     dueDate,
		CreatedBy:   creatorID,
		CreatedAt:   now,
		UpdatedAt:   now,
		TaskNumber:  num,
		TaskKey:     fmt.Sprintf("%s-%d", project.Prefix, num),
	}

	if err := s.taskRepo.Create(ctx, t); err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}
	return t, nil
}

func (s *TaskService) List(ctx context.Context, filter repository.TaskFilter) ([]domain.Task, int, error) {
	return s.taskRepo.List(ctx, filter)
}

func (s *TaskService) Update(ctx context.Context, taskID, callerID uuid.UUID, title, description *string, status *domain.TaskStatus, priority *domain.TaskPriority, assigneeID *string, startDate *string, dueDate *string, blockedReason *string, blockedByTask *string) (*domain.Task, error) {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if title != nil {
		t.Title = *title
	}
	if description != nil {
		t.Description = *description
	}
	if status != nil {
		if *status != t.Status && s.cfSvc != nil {
			allowed, current, max, err := s.cfSvc.CheckWIPLimit(ctx, t.ProjectID, *status)
			if err != nil {
				return nil, fmt.Errorf("checking WIP limit: %w", err)
			}
			if !allowed {
				return nil, fmt.Errorf("WIP limit reached for %s (%d/%d): %w", *status, current, max, domain.ErrConflict)
			}
		}
		t.Status = *status
		if *status == domain.StatusBlocked {
			if blockedReason != nil {
				t.BlockedReason = *blockedReason
			}
			if blockedByTask != nil {
				t.BlockedByTask = *blockedByTask
			}
		} else {
			t.BlockedReason = ""
			t.BlockedByTask = ""
		}
	}
	if priority != nil {
		t.Priority = *priority
	}
	if assigneeID != nil {
		if *assigneeID == "" {
			t.AssigneeID = nil
		} else {
			uid, err := uuid.Parse(*assigneeID)
			if err != nil {
				return nil, domain.ErrInvalidInput
			}
			t.AssigneeID = &uid
		}
	}
	if startDate != nil {
		if *startDate == "" {
			t.StartDate = nil
		} else {
			d, err := time.Parse("2006-01-02", *startDate)
			if err != nil {
				return nil, domain.ErrInvalidInput
			}
			t.StartDate = &d
		}
	}
	if dueDate != nil {
		if *dueDate == "" {
			t.DueDate = nil
		} else {
			d, err := time.Parse("2006-01-02", *dueDate)
			if err != nil {
				return nil, domain.ErrInvalidInput
			}
			t.DueDate = &d
		}
	}
	t.UpdatedAt = time.Now().UTC()

	if err := s.taskRepo.Update(ctx, t); err != nil {
		return nil, fmt.Errorf("updating task: %w", err)
	}
	return t, nil
}

func (s *TaskService) Delete(ctx context.Context, taskID, callerID uuid.UUID) error {
	t, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	p, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil {
		return err
	}

	if t.CreatedBy != callerID && p.OwnerID != callerID {
		return domain.ErrForbidden
	}

	return s.taskRepo.Delete(ctx, taskID)
}

func (s *TaskService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}
