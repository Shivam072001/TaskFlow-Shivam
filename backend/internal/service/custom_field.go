package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shivam/taskflow/backend/internal/domain"
	"github.com/shivam/taskflow/backend/internal/repository"
)

type CustomFieldService struct {
	cfRepo *repository.CustomFieldRepository
}

func NewCustomFieldService(cfr *repository.CustomFieldRepository) *CustomFieldService {
	return &CustomFieldService{cfRepo: cfr}
}

func (s *CustomFieldService) CreateDefinition(ctx context.Context, projectID uuid.UUID, name string, fieldType domain.CustomFieldType, options []string, required bool, createdBy uuid.UUID, callerRole domain.WorkspaceRole) (*domain.CustomFieldDefinition, error) {
	if !callerRole.CanManageCustomFields() {
		return nil, domain.ErrForbidden
	}

	def := &domain.CustomFieldDefinition{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      name,
		FieldType: fieldType,
		Options:   options,
		Required:  required,
		CreatedBy: createdBy,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.cfRepo.CreateDefinition(ctx, def); err != nil {
		return nil, fmt.Errorf("creating custom field definition: %w", err)
	}
	return def, nil
}

func (s *CustomFieldService) ListDefinitions(ctx context.Context, projectID uuid.UUID) ([]domain.CustomFieldDefinition, error) {
	return s.cfRepo.ListDefinitions(ctx, projectID)
}

func (s *CustomFieldService) DeleteDefinition(ctx context.Context, defID uuid.UUID, callerRole domain.WorkspaceRole) error {
	if !callerRole.CanManageCustomFields() {
		return domain.ErrForbidden
	}
	return s.cfRepo.DeleteDefinition(ctx, defID)
}

func (s *CustomFieldService) SetFieldValue(ctx context.Context, taskID, fieldID uuid.UUID, value string) (*domain.CustomFieldValue, error) {
	_, err := s.cfRepo.GetDefinitionByID(ctx, fieldID)
	if err != nil {
		return nil, fmt.Errorf("field definition lookup: %w", err)
	}

	val := &domain.CustomFieldValue{
		ID:      uuid.New(),
		TaskID:  taskID,
		FieldID: fieldID,
		Value:   value,
	}

	if err := s.cfRepo.UpsertValue(ctx, val); err != nil {
		return nil, fmt.Errorf("setting field value: %w", err)
	}
	return val, nil
}

func (s *CustomFieldService) GetFieldValues(ctx context.Context, taskID uuid.UUID) ([]domain.CustomFieldValue, error) {
	return s.cfRepo.GetValues(ctx, taskID)
}

func (s *CustomFieldService) GetWIPLimits(ctx context.Context, projectID uuid.UUID) ([]domain.WIPLimit, error) {
	return s.cfRepo.GetWIPLimits(ctx, projectID)
}

func (s *CustomFieldService) SetWIPLimit(ctx context.Context, projectID uuid.UUID, status domain.TaskStatus, maxTasks int, callerRole domain.WorkspaceRole) (*domain.WIPLimit, error) {
	if !callerRole.CanSetWIPLimits() {
		return nil, domain.ErrForbidden
	}

	limit := &domain.WIPLimit{
		ID:        uuid.New(),
		ProjectID: projectID,
		Status:    status,
		MaxTasks:  maxTasks,
	}

	if err := s.cfRepo.UpsertWIPLimit(ctx, limit); err != nil {
		return nil, fmt.Errorf("setting WIP limit: %w", err)
	}
	return limit, nil
}

func (s *CustomFieldService) DeleteWIPLimit(ctx context.Context, projectID uuid.UUID, status domain.TaskStatus, callerRole domain.WorkspaceRole) error {
	if !callerRole.CanSetWIPLimits() {
		return domain.ErrForbidden
	}
	return s.cfRepo.DeleteWIPLimit(ctx, projectID, status)
}

func (s *CustomFieldService) CheckWIPLimit(ctx context.Context, projectID uuid.UUID, status domain.TaskStatus) (allowed bool, current int, max int, err error) {
	limits, err := s.cfRepo.GetWIPLimits(ctx, projectID)
	if err != nil {
		return false, 0, 0, err
	}

	var limit *domain.WIPLimit
	for _, l := range limits {
		if l.Status == status {
			limit = &l
			break
		}
	}
	if limit == nil {
		return true, 0, 0, nil
	}

	current, err = s.cfRepo.CountTasksByStatus(ctx, projectID, status)
	if err != nil {
		return false, 0, 0, err
	}

	return current < limit.MaxTasks, current, limit.MaxTasks, nil
}
