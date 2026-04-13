import { useQuery } from '@tanstack/react-query';
import * as api from '@core/api/orgDashboard';

export function useOrgMemberStats(orgId: string) {
  return useQuery({
    queryKey: ['org-dashboard', orgId, 'member-stats'],
    queryFn: () => api.fetchOrgMemberStats(orgId),
    enabled: !!orgId,
  });
}

export function useOrgMemberTasks(
  orgId: string,
  userId: string,
  filters?: { status?: string; priority?: string; search?: string; page?: number; limit?: number },
) {
  return useQuery({
    queryKey: ['org-dashboard', orgId, 'member-tasks', userId, filters],
    queryFn: () => api.fetchOrgMemberTasks(orgId, userId, filters),
    enabled: !!orgId && !!userId,
  });
}
