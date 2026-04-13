import { useState, useCallback } from 'react';
import { MessageSquare, Reply, Pencil, Trash2, Loader2, Send } from 'lucide-react';
import type { Comment } from '@/types';
import { useComments, useCreateComment, useUpdateComment, useDeleteComment } from '@hooks/useComments';
import { useConfirm } from '@hooks/useConfirm';
import { RichTextDisplay } from '@components/ui/RichTextDisplay';
import { RichTextEditor } from '@components/ui/RichTextEditor';
import { formatRelative } from '@utils/format';
import { cn } from '@utils/cn';

interface Props {
  entityType: 'project' | 'task';
  entityId: string;
  currentUserId: string;
}

function UserAvatar({ name, className }: { name: string; className?: string }) {
  return (
    <div
      className={cn(
        'flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-sm font-medium text-primary',
        className,
      )}
    >
      {name.charAt(0).toUpperCase()}
    </div>
  );
}

function CommentItem({
  comment,
  currentUserId,
  onReply,
  onUpdate,
  onDelete,
  isUpdating: parentUpdating,
  isDeleting: parentDeleting,
  depth,
}: {
  comment: Comment;
  currentUserId: string;
  onReply: (parentId: string, content: string) => void;
  onUpdate: (commentId: string, content: string) => void;
  onDelete: (commentId: string) => void;
  isUpdating: boolean;
  isDeleting: boolean;
  depth: number;
}) {
  const [replyOpen, setReplyOpen] = useState(false);
  const [replyContent, setReplyContent] = useState('');
  const [editing, setEditing] = useState(false);
  const [editContent, setEditContent] = useState(comment.content);
  const isOwn = comment.user_id === currentUserId;
  const confirm = useConfirm();

  function handleSubmitReply() {
    if (!replyContent.trim() || replyContent === '<p></p>') return;
    onReply(comment.id, replyContent);
    setReplyContent('');
    setReplyOpen(false);
  }

  function handleSubmitEdit() {
    if (!editContent.trim() || editContent === '<p></p>') return;
    onUpdate(comment.id, editContent);
    setEditing(false);
  }

  async function handleDelete() {
    const ok = await confirm({
      title: 'Delete comment',
      description: 'This cannot be undone. All replies will also be removed.',
      confirmLabel: 'Delete',
      variant: 'danger',
    });
    if (ok) onDelete(comment.id);
  }

  return (
    <div className={cn(depth > 0 && 'ml-8 border-l-2 border-border pl-4')}>
      <div className="group flex gap-3 py-3">
        <UserAvatar name={comment.user_name} />
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-foreground">{comment.user_name}</span>
            <span className="text-xs text-muted-foreground">{formatRelative(comment.created_at)}</span>
            {comment.updated_at !== comment.created_at && (
              <span className="text-xs text-muted-foreground italic">(edited)</span>
            )}
          </div>

          {editing ? (
            <div className="mt-2 space-y-2">
              <RichTextEditor value={editContent} onChange={setEditContent} placeholder="Edit comment…" />
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={handleSubmitEdit}
                  disabled={parentUpdating}
                  className="inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                >
                  {parentUpdating ? <Loader2 className="h-3 w-3 animate-spin" /> : <Send className="h-3 w-3" />}
                  Save
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setEditing(false);
                    setEditContent(comment.content);
                  }}
                  className="rounded-md px-3 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted"
                >
                  Cancel
                </button>
              </div>
            </div>
          ) : (
            <div className="mt-1">
              <RichTextDisplay content={comment.content} />
            </div>
          )}

          {!editing && (
            <div className="mt-1.5 flex items-center gap-3 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100">
              {depth === 0 && (
                <button
                  type="button"
                  onClick={() => setReplyOpen((o) => !o)}
                  className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
                >
                  <Reply className="h-3.5 w-3.5" />
                  Reply
                </button>
              )}
              {isOwn && (
                <>
                  <button
                    type="button"
                    onClick={() => {
                      setEditContent(comment.content);
                      setEditing(true);
                    }}
                    className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
                  >
                    <Pencil className="h-3.5 w-3.5" />
                    Edit
                  </button>
                  <button
                    type="button"
                    onClick={handleDelete}
                    disabled={parentDeleting}
                    className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-destructive disabled:opacity-50"
                  >
                    {parentDeleting ? (
                      <Loader2 className="h-3.5 w-3.5 animate-spin" />
                    ) : (
                      <Trash2 className="h-3.5 w-3.5" />
                    )}
                    Delete
                  </button>
                </>
              )}
            </div>
          )}
        </div>
      </div>

      {replyOpen && (
        <div className="mb-2 ml-11 space-y-2">
          <RichTextEditor value={replyContent} onChange={setReplyContent} placeholder="Write a reply…" />
          <div className="flex gap-2">
            <button
              type="button"
              onClick={handleSubmitReply}
              className="inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground hover:bg-primary/90"
            >
              <Send className="h-3 w-3" />
              Reply
            </button>
            <button
              type="button"
              onClick={() => {
                setReplyOpen(false);
                setReplyContent('');
              }}
              className="rounded-md px-3 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted"
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {comment.replies?.map((reply) => (
        <CommentItem
          key={reply.id}
          comment={reply}
          currentUserId={currentUserId}
          onReply={onReply}
          onUpdate={onUpdate}
          onDelete={onDelete}
          isUpdating={parentUpdating}
          isDeleting={parentDeleting}
          depth={depth + 1}
        />
      ))}
    </div>
  );
}

export function CommentThread({ entityType, entityId, currentUserId }: Props) {
  const { data: comments, isLoading } = useComments(entityType, entityId);
  const createComment = useCreateComment(entityType, entityId);
  const updateComment = useUpdateComment(entityType, entityId);
  const deleteComment = useDeleteComment(entityType, entityId);

  const [newContent, setNewContent] = useState('');

  const handleReply = useCallback(
    (parentId: string, content: string) => createComment.mutate({ content, parentId }),
    [createComment],
  );

  const handleUpdate = useCallback(
    (commentId: string, content: string) => updateComment.mutate({ commentId, content }),
    [updateComment],
  );

  const handleDeleteComment = useCallback(
    (commentId: string) => deleteComment.mutate(commentId),
    [deleteComment],
  );

  const handleCreate = useCallback(() => {
    if (!newContent.trim() || newContent === '<p></p>') return;
    createComment.mutate({ content: newContent }, { onSuccess: () => setNewContent('') });
  }, [newContent, createComment]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <Loader2 className="h-5 w-5 animate-spin" />
      </div>
    );
  }

  const list = comments ?? [];

  return (
    <div className="space-y-1">
      <h3 className="flex items-center gap-2 text-sm font-semibold text-foreground">
        <MessageSquare className="h-4 w-4" />
        Comments
        {list.length > 0 && (
          <span className="rounded-full bg-muted px-2 py-0.5 text-xs font-normal text-muted-foreground">
            {list.length}
          </span>
        )}
      </h3>

      {list.length === 0 ? (
        <p className="py-6 text-center text-sm text-muted-foreground">No comments yet. Start the conversation!</p>
      ) : (
        <div className="divide-y divide-border">
          {list.map((comment) => (
            <CommentItem
              key={comment.id}
              comment={comment}
              currentUserId={currentUserId}
              onReply={handleReply}
              onUpdate={handleUpdate}
              onDelete={handleDeleteComment}
              isUpdating={updateComment.isPending && updateComment.variables?.commentId === comment.id}
              isDeleting={deleteComment.isPending && deleteComment.variables === comment.id}
              depth={0}
            />
          ))}
        </div>
      )}

      <div className="space-y-2 pt-4">
        <RichTextEditor value={newContent} onChange={setNewContent} placeholder="Add a comment…" />
        <button
          type="button"
          onClick={handleCreate}
          disabled={createComment.isPending || !newContent.trim() || newContent === '<p></p>'}
          className="inline-flex items-center gap-1.5 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
        >
          {createComment.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
          Post Comment
        </button>
      </div>
    </div>
  );
}
