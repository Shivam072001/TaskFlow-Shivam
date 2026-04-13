import type { Comment } from '@/types';
import apiClient from './client';

export async function listComments(entityType: string, entityId: string): Promise<Comment[]> {
  const { data } = await apiClient.get<{ comments: Comment[] }>(
    `/${entityType}s/${entityId}/comments`,
  );
  return data.comments;
}

export async function createComment(
  entityType: string,
  entityId: string,
  content: string,
  parentId?: string,
): Promise<Comment> {
  const { data } = await apiClient.post<Comment>(`/${entityType}s/${entityId}/comments`, {
    content,
    parent_id: parentId ?? null,
  });
  return data;
}

export async function updateComment(commentId: string, content: string): Promise<Comment> {
  const { data } = await apiClient.patch<Comment>(`/comments/${commentId}`, { content });
  return data;
}

export async function deleteComment(commentId: string): Promise<void> {
  await apiClient.delete(`/comments/${commentId}`);
}
