import type { Task, PaginatedResponse } from '@/types';
import apiClient from './client';

export interface TaskFilters {
  status?: string;
  assignee?: string;
  priority?: string;
  search?: string;
  start_date?: string;
  page?: number;
  limit?: number;
}

export async function listTasks(projectId: string, filters?: TaskFilters): Promise<PaginatedResponse<Task>> {
  const params = new URLSearchParams();
  if (filters?.status) params.set('status', filters.status);
  if (filters?.assignee) params.set('assignee', filters.assignee);
  if (filters?.priority) params.set('priority', filters.priority);
  if (filters?.search) params.set('search', filters.search);
  if (filters?.start_date) params.set('start_date', filters.start_date);
  if (filters?.page) params.set('page', String(filters.page));
  if (filters?.limit) params.set('limit', String(filters.limit));

  const { data } = await apiClient.get<PaginatedResponse<Task>>(
    `/projects/${projectId}/tasks?${params.toString()}`
  );
  return data;
}

export interface CreateTaskBody {
  title: string;
  description?: string;
  priority?: string;
  assignee_id?: string;
  start_date?: string;
  due_date?: string;
}

export async function createTask(projectId: string, body: CreateTaskBody): Promise<Task> {
  const { data } = await apiClient.post<Task>(`/projects/${projectId}/tasks`, body);
  return data;
}

export interface UpdateTaskBody {
  title?: string;
  description?: string;
  status?: string;
  priority?: string;
  assignee_id?: string;
  start_date?: string;
  due_date?: string;
  blocked_reason?: string;
  blocked_by_task?: string;
}

export async function updateTask(taskId: string, body: UpdateTaskBody): Promise<Task> {
  const { data } = await apiClient.patch<Task>(`/tasks/${taskId}`, body);
  return data;
}

export async function deleteTask(taskId: string): Promise<void> {
  await apiClient.delete(`/tasks/${taskId}`);
}
