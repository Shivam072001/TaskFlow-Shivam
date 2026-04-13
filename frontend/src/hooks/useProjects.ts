import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { ProjectListParams } from '@core/api/projects';
import * as api from '@core/api/projects';

export function useProjects(orgId: string, workspaceId: string, params: ProjectListParams = {}) {
  return useQuery({
    queryKey: ['projects', orgId, workspaceId, params],
    queryFn: () => api.listProjects(orgId, workspaceId, params),
    enabled: !!orgId && !!workspaceId,
  });
}

export function useProject(id: string) {
  return useQuery({
    queryKey: ['project', id],
    queryFn: () => api.getProject(id),
    enabled: !!id,
  });
}

export function useCreateProject(orgId: string, workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ name, prefix, description }: { name: string; prefix: string; description: string }) =>
      api.createProject(orgId, workspaceId, name, prefix, description),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['projects', orgId, workspaceId] });
      qc.invalidateQueries({ queryKey: ['workspace-stats', workspaceId] });
      qc.invalidateQueries({ queryKey: ['org-prefixes', orgId] });
    },
  });
}

export function useDeleteProject(orgId: string, workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deleteProject(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['projects', orgId, workspaceId] });
      qc.invalidateQueries({ queryKey: ['workspace-stats', workspaceId] });
      qc.invalidateQueries({ queryKey: ['org-prefixes', orgId] });
    },
  });
}
