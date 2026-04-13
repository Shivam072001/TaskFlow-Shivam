import type { WorkspaceInvitation, WorkspaceRole } from '@/types';
import apiClient from './client';

export async function sendInvite(
  orgId: string,
  workspaceId: string,
  email: string,
  role: WorkspaceRole,
): Promise<WorkspaceInvitation> {
  const { data } = await apiClient.post<WorkspaceInvitation>(
    `/organizations/${orgId}/workspaces/${workspaceId}/members`,
    { email, role },
  );
  return data;
}

export async function listMyInvitations(): Promise<WorkspaceInvitation[]> {
  const { data } = await apiClient.get<{ invitations: WorkspaceInvitation[] }>('/invitations');
  return data.invitations;
}

export async function listWorkspaceInvitations(orgId: string, workspaceId: string): Promise<WorkspaceInvitation[]> {
  const { data } = await apiClient.get<{ invitations: WorkspaceInvitation[] }>(
    `/organizations/${orgId}/workspaces/${workspaceId}/invitations`,
  );
  return data.invitations;
}

export async function respondToInvitation(invitationId: string, accept: boolean): Promise<void> {
  await apiClient.patch(`/invitations/${invitationId}`, {
    action: accept ? 'accept' : 'decline',
  });
}
