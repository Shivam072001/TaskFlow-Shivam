import type { OrgInvitation, OrgRole } from '@/types';
import apiClient from './client';

export async function sendOrgInvite(orgId: string, email: string, role: OrgRole): Promise<OrgInvitation> {
  const { data } = await apiClient.post<OrgInvitation>(
    `/organizations/${orgId}/invitations`,
    { email, role },
  );
  return data;
}

export async function listOrgInvitations(orgId: string): Promise<OrgInvitation[]> {
  const { data } = await apiClient.get<{ invitations: OrgInvitation[] }>(
    `/organizations/${orgId}/invitations`,
  );
  return data.invitations;
}

export async function listMyOrgInvitations(): Promise<OrgInvitation[]> {
  const { data } = await apiClient.get<{ invitations: OrgInvitation[] }>('/org-invitations');
  return data.invitations;
}

export async function respondToOrgInvitation(invitationId: string, accept: boolean): Promise<void> {
  await apiClient.patch(`/org-invitations/${invitationId}`, {
    action: accept ? 'accept' : 'decline',
  });
}
