import type { Task, PaginatedResponse } from '@/types';
import apiClient from './client';

export interface MyTaskFilters {
  status?: string;
  priority?: string;
  search?: string;
  due_before?: string;
  due_after?: string;
  project_id?: string;
  page?: number;
  limit?: number;
}

export interface UserTaskStats {
  total: number;
  by_status: Record<string, number>;
  by_priority: Record<string, number>;
  overdue: number;
  due_today: number;
  due_this_week: number;
  completed: number;
}

export interface ProjectOption {
  id: string;
  name: string;
}

export async function fetchMyTasks(filters?: MyTaskFilters): Promise<PaginatedResponse<Task>> {
  const params = new URLSearchParams();
  if (filters?.status) params.set('status', filters.status);
  if (filters?.priority) params.set('priority', filters.priority);
  if (filters?.search) params.set('search', filters.search);
  if (filters?.due_before) params.set('due_before', filters.due_before);
  if (filters?.due_after) params.set('due_after', filters.due_after);
  if (filters?.project_id) params.set('project_id', filters.project_id);
  if (filters?.page) params.set('page', String(filters.page));
  if (filters?.limit) params.set('limit', String(filters.limit));
  const { data } = await apiClient.get<PaginatedResponse<Task>>(`/dashboard/my-tasks?${params}`);
  return data;
}

export async function fetchMyStats(): Promise<UserTaskStats> {
  const { data } = await apiClient.get<UserTaskStats>('/dashboard/my-stats');
  return data;
}

export async function fetchMyProjects(): Promise<ProjectOption[]> {
  const { data } = await apiClient.get<ProjectOption[]>('/dashboard/my-projects');
  return data;
}
