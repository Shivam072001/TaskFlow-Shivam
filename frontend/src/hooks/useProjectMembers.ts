import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/projectMembers';
import type { WorkspaceRole } from '@/types';

export function useProjectMembers(orgId: string, workspaceId: string, projectId: string) {
  return useQuery({
    queryKey: ['project-members', projectId],
    queryFn: () => api.listProjectMembers(orgId, workspaceId, projectId),
    enabled: !!orgId && !!workspaceId && !!projectId,
  });
}

export function useAddProjectMember(orgId: string, workspaceId: string, projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: WorkspaceRole }) =>
      api.addProjectMember(orgId, workspaceId, projectId, userId, role),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['project-members', projectId] }),
  });
}

export function useRemoveProjectMember(orgId: string, workspaceId: string, projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) =>
      api.removeProjectMember(orgId, workspaceId, projectId, userId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['project-members', projectId] }),
  });
}

export function useUpdateProjectMemberRole(orgId: string, workspaceId: string, projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: WorkspaceRole }) =>
      api.updateProjectMemberRole(orgId, workspaceId, projectId, userId, role),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['project-members', projectId] }),
  });
}

export function useLeaveProject(orgId: string, workspaceId: string, projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api.leaveProject(orgId, workspaceId, projectId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['project-members', projectId] }),
  });
}
