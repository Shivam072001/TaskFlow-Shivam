import type { ProjectMember, WorkspaceRole } from '@/types';
import apiClient from './client';

export async function listProjectMembers(
  orgId: string, workspaceId: string, projectId: string,
): Promise<ProjectMember[]> {
  const { data } = await apiClient.get<{ members: ProjectMember[] }>(
    `/organizations/${orgId}/workspaces/${workspaceId}/projects/${projectId}/members`,
  );
  return data.members;
}

export async function addProjectMember(
  orgId: string, workspaceId: string, projectId: string, userId: string, role: WorkspaceRole,
): Promise<ProjectMember> {
  const { data } = await apiClient.post<ProjectMember>(
    `/organizations/${orgId}/workspaces/${workspaceId}/projects/${projectId}/members`,
    { user_id: userId, role },
  );
  return data;
}

export async function removeProjectMember(
  orgId: string, workspaceId: string, projectId: string, userId: string,
): Promise<void> {
  await apiClient.delete(
    `/organizations/${orgId}/workspaces/${workspaceId}/projects/${projectId}/members/${userId}`,
  );
}

export async function updateProjectMemberRole(
  orgId: string, workspaceId: string, projectId: string, userId: string, role: WorkspaceRole,
): Promise<void> {
  await apiClient.patch(
    `/organizations/${orgId}/workspaces/${workspaceId}/projects/${projectId}/members/${userId}`,
    { role },
  );
}

export async function leaveProject(
  orgId: string, workspaceId: string, projectId: string,
): Promise<void> {
  await apiClient.post(
    `/organizations/${orgId}/workspaces/${workspaceId}/projects/${projectId}/members/leave`,
  );
}
