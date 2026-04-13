import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/organizations';

export function useOrganizations() {
  return useQuery({
    queryKey: ['organizations'],
    queryFn: api.listOrganizations,
  });
}

export function useOrganization(orgId: string) {
  return useQuery({
    queryKey: ['organization', orgId],
    queryFn: () => api.getOrganization(orgId),
    enabled: !!orgId,
  });
}

export function useCreateOrganization() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ name, slug }: { name: string; slug: string }) =>
      api.createOrganization(name, slug),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['organizations'] }),
  });
}

export function useOrgMembers(orgId: string) {
  return useQuery({
    queryKey: ['org-members', orgId],
    queryFn: () => api.listOrgMembers(orgId),
    enabled: !!orgId,
  });
}

export function useOrgPrefixes(orgId: string) {
  return useQuery({
    queryKey: ['org-prefixes', orgId],
    queryFn: () => api.listOrgPrefixes(orgId),
    enabled: !!orgId,
  });
}

export function useTaskByKey(orgId: string, taskKey: string) {
  return useQuery({
    queryKey: ['task-by-key', orgId, taskKey],
    queryFn: () => api.getTaskByKey(orgId, taskKey),
    enabled: !!orgId && !!taskKey,
  });
}
