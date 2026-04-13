import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/workspaces';

export function useWorkspaces(orgId: string, params: { page?: number; limit?: number } = {}) {
  return useQuery({
    queryKey: ['workspaces', orgId, params],
    queryFn: () => api.listWorkspaces(orgId, params),
    enabled: !!orgId,
  });
}

export function useWorkspace(orgId: string, id: string) {
  return useQuery({
    queryKey: ['workspace', orgId, id],
    queryFn: () => api.getWorkspace(orgId, id),
    enabled: !!orgId && !!id,
  });
}

export function useWorkspaceStats(orgId: string, id: string) {
  return useQuery({
    queryKey: ['workspace-stats', id],
    queryFn: () => api.getWorkspaceStats(orgId, id),
    enabled: !!orgId && !!id,
  });
}

export function useWorkspaceMembers(orgId: string, id: string) {
  return useQuery({
    queryKey: ['workspace-members', id],
    queryFn: () => api.listMembers(orgId, id),
    enabled: !!orgId && !!id,
  });
}

export function useCreateWorkspace(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ name, description }: { name: string; description: string }) =>
      api.createWorkspace(orgId, name, description),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['workspaces', orgId] }),
  });
}

export function useDeleteWorkspace(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deleteWorkspace(orgId, id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['workspaces', orgId] }),
  });
}

export function useUpdateMemberRole(orgId: string, workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: string }) =>
      api.updateMemberRole(orgId, workspaceId, userId, role),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['workspace-members', workspaceId] }),
  });
}

export function useRemoveMember(orgId: string, workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) => api.removeMember(orgId, workspaceId, userId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['workspace-members', workspaceId] }),
  });
}

export function useDirectAddMember(orgId: string, workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: string }) =>
      api.directAddMember(orgId, workspaceId, userId, role),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['workspace-members', workspaceId] }),
  });
}

export function useLeaveWorkspace(orgId: string, workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api.leaveWorkspace(orgId, workspaceId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['workspace-members', workspaceId] });
      qc.invalidateQueries({ queryKey: ['workspaces', orgId] });
    },
  });
}

export function useLeaveOrg(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api.leaveOrg(orgId),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['organizations'] });
    },
  });
}

export function useRemoveOrgMember(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) => api.removeOrgMember(orgId, userId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['org-members', orgId] }),
  });
}

export function useUpdateOrgMemberRole(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: string }) =>
      api.updateOrgMemberRole(orgId, userId, role),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['org-members', orgId] }),
  });
}
