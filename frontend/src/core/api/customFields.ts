import type { CustomFieldDefinition, CustomFieldValue, WIPLimit, TaskStatus } from '@/types';
import apiClient from './client';

export async function listDefinitions(projectId: string): Promise<CustomFieldDefinition[]> {
  const { data } = await apiClient.get<{ custom_fields: CustomFieldDefinition[] }>(
    `/projects/${projectId}/custom-fields`,
  );
  return data.custom_fields;
}

export interface CreateDefinitionBody {
  name: string;
  field_type: 'text' | 'number' | 'select';
  options?: string[];
  required?: boolean;
}

export async function createDefinition(
  projectId: string,
  body: CreateDefinitionBody,
): Promise<CustomFieldDefinition> {
  const { data } = await apiClient.post<CustomFieldDefinition>(
    `/projects/${projectId}/custom-fields`,
    body,
  );
  return data;
}

export async function deleteDefinition(projectId: string, defId: string): Promise<void> {
  await apiClient.delete(`/projects/${projectId}/custom-fields/${defId}`);
}

export async function getFieldValues(taskId: string): Promise<CustomFieldValue[]> {
  const { data } = await apiClient.get<{ custom_fields: CustomFieldValue[] }>(
    `/tasks/${taskId}/custom-fields`,
  );
  return data.custom_fields;
}

export async function setFieldValue(taskId: string, fieldId: string, value: string): Promise<void> {
  await apiClient.put(`/tasks/${taskId}/custom-fields/${fieldId}`, { value });
}

export async function getWIPLimits(projectId: string): Promise<WIPLimit[]> {
  const { data } = await apiClient.get<{ wip_limits: WIPLimit[] }>(
    `/projects/${projectId}/wip-limits`,
  );
  return data.wip_limits;
}

export async function setWIPLimit(
  projectId: string,
  status: TaskStatus,
  maxTasks: number,
): Promise<void> {
  await apiClient.put(`/projects/${projectId}/wip-limits`, { status, max_tasks: maxTasks });
}

export async function deleteWIPLimit(projectId: string, status: TaskStatus): Promise<void> {
  await apiClient.delete(`/projects/${projectId}/wip-limits?status=${status}`);
}
