import type { Project, PaginatedResponse } from '@/types';
import apiClient from './client';

export interface ProjectListParams {
  page?: number;
  limit?: number;
  search?: string;
  owner?: string;
}

export async function listProjects(orgId: string, workspaceId: string, params: ProjectListParams = {}): Promise<PaginatedResponse<Project>> {
  const { data } = await apiClient.get<PaginatedResponse<Project>>(`/organizations/${orgId}/workspaces/${workspaceId}/projects`, { params });
  return data;
}

export async function getProject(id: string): Promise<Project> {
  const { data } = await apiClient.get<Project>(`/projects/${id}`);
  return data;
}

export async function createProject(orgId: string, workspaceId: string, name: string, prefix: string, description: string): Promise<Project> {
  const { data } = await apiClient.post<Project>(`/organizations/${orgId}/workspaces/${workspaceId}/projects`, { name, prefix, description });
  return data;
}

export async function updateProject(id: string, body: { name?: string; description?: string }): Promise<Project> {
  const { data } = await apiClient.patch<Project>(`/projects/${id}`, body);
  return data;
}

export async function deleteProject(id: string): Promise<void> {
  await apiClient.delete(`/projects/${id}`);
}
