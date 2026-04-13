import type { Task, PaginatedResponse } from '@/types';
import apiClient from './client';

export interface OrgMemberStats {
  user_id: string;
  user_name: string;
  user_email: string;
  role: string;
  total: number;
  by_status: Record<string, number>;
  overdue: number;
  completed: number;
}

export async function fetchOrgMemberStats(orgId: string): Promise<OrgMemberStats[]> {
  const { data } = await apiClient.get<OrgMemberStats[]>(`/organizations/${orgId}/dashboard/member-stats`);
  return data;
}

export async function fetchOrgMemberTasks(
  orgId: string,
  userId: string,
  filters?: { status?: string; priority?: string; search?: string; page?: number; limit?: number },
): Promise<PaginatedResponse<Task>> {
  const params = new URLSearchParams();
  if (filters?.status) params.set('status', filters.status);
  if (filters?.priority) params.set('priority', filters.priority);
  if (filters?.search) params.set('search', filters.search);
  if (filters?.page) params.set('page', String(filters.page));
  if (filters?.limit) params.set('limit', String(filters.limit));
  const { data } = await apiClient.get<PaginatedResponse<Task>>(
    `/organizations/${orgId}/dashboard/member-tasks/${userId}?${params}`,
  );
  return data;
}
