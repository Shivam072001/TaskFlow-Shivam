import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/tasks';
import type { Task } from '@/types';

interface TaskFilters {
  status?: string;
  assignee?: string;
  priority?: string;
  search?: string;
}

export function useTasks(projectId: string, filters?: TaskFilters) {
  return useQuery({
    queryKey: ['tasks', projectId, filters],
    queryFn: () => api.listTasks(projectId, { ...filters, limit: 200 }),
    enabled: !!projectId,
  });
}

export function useCreateTask(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: api.CreateTaskBody) => api.createTask(projectId, body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tasks', projectId] }),
  });
}

export function useUpdateTask(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ taskId, body }: { taskId: string; body: api.UpdateTaskBody }) =>
      api.updateTask(taskId, body),
    onMutate: async ({ taskId, body }) => {
      await qc.cancelQueries({ queryKey: ['tasks', projectId] });
      const prev = qc.getQueriesData<{ data: Task[]; meta: unknown }>({ queryKey: ['tasks', projectId] });

      qc.setQueriesData<{ data: Task[]; meta: unknown }>(
        { queryKey: ['tasks', projectId] },
        (old) => {
          if (!old) return old;
          return {
            ...old,
            data: old.data.map((t: Task) =>
              t.id === taskId ? { ...t, ...body, updated_at: new Date().toISOString() } : t
            ),
          };
        },
      );

      return { prev };
    },
    onError: (_err, _vars, context) => {
      if (context?.prev) {
        context.prev.forEach(([key, data]) => qc.setQueryData(key, data));
      }
    },
    onSettled: () => qc.invalidateQueries({ queryKey: ['tasks', projectId] }),
  });
}

export function useDeleteTask(projectId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (taskId: string) => api.deleteTask(taskId),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['tasks', projectId] }),
  });
}
