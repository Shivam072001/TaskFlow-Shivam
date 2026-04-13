import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import * as api from '@core/api/comments';

export function useComments(entityType: string, entityId: string) {
  return useQuery({
    queryKey: ['comments', entityType, entityId],
    queryFn: () => api.listComments(entityType, entityId),
    enabled: !!entityId,
  });
}

export function useCreateComment(entityType: string, entityId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ content, parentId }: { content: string; parentId?: string }) =>
      api.createComment(entityType, entityId, content, parentId),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['comments', entityType, entityId] }),
  });
}

export function useUpdateComment(entityType: string, entityId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ commentId, content }: { commentId: string; content: string }) =>
      api.updateComment(commentId, content),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['comments', entityType, entityId] }),
  });
}

export function useDeleteComment(entityType: string, entityId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (commentId: string) => api.deleteComment(commentId),
    onSuccess: () =>
      qc.invalidateQueries({ queryKey: ['comments', entityType, entityId] }),
  });
}
