import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/teams';

export function useTeams(orgId: string) {
  return useQuery({
    queryKey: ['teams', orgId],
    queryFn: () => api.listTeams(orgId),
    enabled: !!orgId,
  });
}

export function useCreateTeam(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (name: string) => api.createTeam(orgId, name),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['teams', orgId] }),
  });
}

export function useTeamDetail(orgId: string, teamId: string) {
  return useQuery({
    queryKey: ['team', teamId],
    queryFn: () => api.getTeam(orgId, teamId),
    enabled: !!orgId && !!teamId,
  });
}

export function useUpdateTeam(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ teamId, name }: { teamId: string; name: string }) =>
      api.updateTeam(orgId, teamId, name),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['teams', orgId] }),
  });
}

export function useDeleteTeam(orgId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (teamId: string) => api.deleteTeam(orgId, teamId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['teams', orgId] }),
  });
}

export function useAddTeamMember(orgId: string, teamId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) => api.addTeamMember(orgId, teamId, userId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['team', teamId] }),
  });
}

export function useRemoveTeamMember(orgId: string, teamId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) => api.removeTeamMember(orgId, teamId, userId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['team', teamId] }),
  });
}

export function useAddTeamToWorkspace(orgId: string, workspaceId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ teamId, defaultRole }: { teamId: string; defaultRole: string }) =>
      api.addTeamToWorkspace(orgId, workspaceId, teamId, defaultRole),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['workspace-members', workspaceId] });
    },
  });
}

export function useAddTeamToProject(orgId: string, workspaceId: string, projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ teamId, defaultRole }: { teamId: string; defaultRole: string }) =>
      api.addTeamToProject(orgId, workspaceId, projectId, teamId, defaultRole),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['project-members', projectId] });
    },
  });
}
