import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/customFields';
import type { TaskStatus } from '@/types';

export function useCustomFieldDefinitions(projectId: string) {
  return useQuery({
    queryKey: ['custom-fields', projectId],
    queryFn: () => api.listDefinitions(projectId),
    enabled: !!projectId,
  });
}

export function useCreateCustomField(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: api.CreateDefinitionBody) => api.createDefinition(projectId, body),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['custom-fields', projectId] }),
  });
}

export function useDeleteCustomField(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (defId: string) => api.deleteDefinition(projectId, defId),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['custom-fields', projectId] }),
  });
}

export function useFieldValues(taskId: string) {
  return useQuery({
    queryKey: ['field-values', taskId],
    queryFn: () => api.getFieldValues(taskId),
    enabled: !!taskId,
  });
}

export function useSetFieldValue(taskId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ fieldId, value }: { fieldId: string; value: string }) =>
      api.setFieldValue(taskId, fieldId, value),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['field-values', taskId] }),
  });
}

export function useWIPLimits(projectId: string) {
  return useQuery({
    queryKey: ['wip-limits', projectId],
    queryFn: () => api.getWIPLimits(projectId),
    enabled: !!projectId,
  });
}

export function useSetWIPLimit(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ status, maxTasks }: { status: TaskStatus; maxTasks: number }) =>
      api.setWIPLimit(projectId, status, maxTasks),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['wip-limits', projectId] }),
  });
}

export function useDeleteWIPLimit(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (status: TaskStatus) => api.deleteWIPLimit(projectId, status),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['wip-limits', projectId] }),
  });
}
