import { useQuery } from '@tanstack/react-query';
import * as api from '@core/api/dashboard';

export function useMyTasks(filters?: api.MyTaskFilters) {
  return useQuery({
    queryKey: ['dashboard', 'my-tasks', filters],
    queryFn: () => api.fetchMyTasks(filters),
  });
}

export function useMyStats() {
  return useQuery({
    queryKey: ['dashboard', 'my-stats'],
    queryFn: api.fetchMyStats,
  });
}

export function useMyProjects() {
  return useQuery({
    queryKey: ['dashboard', 'my-projects'],
    queryFn: api.fetchMyProjects,
  });
}
