import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/orgInvitations';
import type { OrgRole } from '@/types';

export function useOrgInvitations(orgId: string) {
  return useQuery({
    queryKey: ['org-invitations', orgId],
    queryFn: () => api.listOrgInvitations(orgId),
    enabled: !!orgId,
  });
}

export function useMyOrgInvitations() {
  return useQuery({
    queryKey: ['org-invitations', 'mine'],
    queryFn: api.listMyOrgInvitations,
  });
}

export function useSendOrgInvite(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ email, role }: { email: string; role: OrgRole }) =>
      api.sendOrgInvite(orgId, email, role),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['org-invitations', orgId] }),
  });
}

export function useRespondOrgInvite() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ invitationId, accept }: { invitationId: string; accept: boolean }) =>
      api.respondToOrgInvitation(invitationId, accept),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['org-invitations'] });
      qc.invalidateQueries({ queryKey: ['organizations'] });
    },
  });
}
