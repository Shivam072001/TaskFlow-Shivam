import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/invitations';
import type { WorkspaceRole } from '@/types';

export function useMyInvitations() {
  return useQuery({
    queryKey: ['invitations', 'mine'],
    queryFn: api.listMyInvitations,
  });
}

export function useWorkspaceInvitations(orgId: string, workspaceId: string) {
  return useQuery({
    queryKey: ['invitations', 'workspace', workspaceId],
    queryFn: () => api.listWorkspaceInvitations(orgId, workspaceId),
    enabled: !!orgId && !!workspaceId,
  });
}

export function useSendInvite(orgId: string, workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ email, role }: { email: string; role: WorkspaceRole }) =>
      api.sendInvite(orgId, workspaceId, email, role),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['invitations', 'workspace', workspaceId] }),
  });
}

export function useRespondToInvitation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ invitationId, accept }: { invitationId: string; accept: boolean }) =>
      api.respondToInvitation(invitationId, accept),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['invitations'] });
      qc.invalidateQueries({ queryKey: ['workspace-members'] });
    },
  });
}
