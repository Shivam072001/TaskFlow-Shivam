import type { Workspace, WorkspaceMember, WorkspaceStats, PaginatedResponse } from '@/types';
import apiClient from './client';

export async function listWorkspaces(orgId: string, params: { page?: number; limit?: number } = {}): Promise<PaginatedResponse<Workspace>> {
  const { data } = await apiClient.get<PaginatedResponse<Workspace>>(`/organizations/${orgId}/workspaces`, { params });
  return data;
}

export async function getWorkspace(orgId: string, id: string): Promise<Workspace> {
  const { data } = await apiClient.get<Workspace>(`/organizations/${orgId}/workspaces/${id}`);
  return data;
}

export async function createWorkspace(orgId: string, name: string, description: string): Promise<Workspace> {
  const { data } = await apiClient.post<Workspace>(`/organizations/${orgId}/workspaces`, { name, description });
  return data;
}

export async function updateWorkspace(orgId: string, id: string, body: { name?: string; description?: string }): Promise<Workspace> {
  const { data } = await apiClient.patch<Workspace>(`/organizations/${orgId}/workspaces/${id}`, body);
  return data;
}

export async function deleteWorkspace(orgId: string, id: string): Promise<void> {
  await apiClient.delete(`/organizations/${orgId}/workspaces/${id}`);
}

export async function getWorkspaceStats(orgId: string, id: string): Promise<WorkspaceStats> {
  const { data } = await apiClient.get<WorkspaceStats>(`/organizations/${orgId}/workspaces/${id}/stats`);
  return data;
}

export async function listMembers(orgId: string, workspaceId: string): Promise<WorkspaceMember[]> {
  const { data } = await apiClient.get<{ members: WorkspaceMember[] }>(`/organizations/${orgId}/workspaces/${workspaceId}/members`);
  return data.members;
}

export async function inviteMember(orgId: string, workspaceId: string, email: string, role: string): Promise<WorkspaceMember> {
  const { data } = await apiClient.post<WorkspaceMember>(`/organizations/${orgId}/workspaces/${workspaceId}/members`, { email, role });
  return data;
}

export async function updateMemberRole(orgId: string, workspaceId: string, userId: string, role: string): Promise<void> {
  await apiClient.patch(`/organizations/${orgId}/workspaces/${workspaceId}/members/${userId}`, { role });
}

export async function removeMember(orgId: string, workspaceId: string, userId: string): Promise<void> {
  await apiClient.delete(`/organizations/${orgId}/workspaces/${workspaceId}/members/${userId}`);
}

export async function directAddMember(orgId: string, workspaceId: string, userId: string, role: string): Promise<void> {
  await apiClient.post(`/organizations/${orgId}/workspaces/${workspaceId}/members/add`, { user_id: userId, role });
}

export async function leaveWorkspace(orgId: string, workspaceId: string): Promise<void> {
  await apiClient.post(`/organizations/${orgId}/workspaces/${workspaceId}/members/leave`);
}

export async function leaveOrg(orgId: string): Promise<void> {
  await apiClient.post(`/organizations/${orgId}/members/leave`);
}

export async function removeOrgMember(orgId: string, userId: string): Promise<void> {
  await apiClient.delete(`/organizations/${orgId}/members/${userId}`);
}
