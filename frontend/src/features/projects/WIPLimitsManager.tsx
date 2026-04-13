import { useState } from 'react';
import { Save, Trash2, Loader2, Gauge } from 'lucide-react';
import type { TaskStatus, WIPLimit } from '@/types';
import { useWIPLimits, useSetWIPLimit, useDeleteWIPLimit } from '@hooks/useCustomFields';
import { statusLabels } from '@utils/format';
import { cn } from '@utils/cn';

interface Props {
  projectId: string;
  canManage: boolean;
}

const statuses: TaskStatus[] = ['todo', 'in_progress', 'blocked', 'done'];

function LimitRow({
  status,
  limit,
  canManage,
  onSave,
  onRemove,
  isSaving,
  isRemoving,
}: {
  status: TaskStatus;
  limit: WIPLimit | undefined;
  canManage: boolean;
  onSave: (status: TaskStatus, maxTasks: number) => void;
  onRemove: (status: TaskStatus) => void;
  isSaving: boolean;
  isRemoving: boolean;
}) {
  const [value, setValue] = useState(limit?.max_tasks?.toString() ?? '');

  function handleSave() {
    const num = parseInt(value, 10);
    if (isNaN(num) || num < 1) return;
    onSave(status, num);
  }

  return (
    <div className="flex items-center gap-4 rounded-md border border-border bg-card px-4 py-3">
      <span className="w-28 text-sm font-medium text-foreground">{statusLabels[status]}</span>

      <span className="text-sm text-muted-foreground">
        {limit ? `Limit: ${limit.max_tasks}` : 'No limit'}
      </span>

      {canManage && (
        <div className="ml-auto flex items-center gap-2">
          <input
            type="number"
            min={1}
            value={value}
            onChange={(e) => setValue(e.target.value)}
            placeholder="Max"
            className="w-20 rounded-md border border-border bg-background px-2 py-1.5 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          />
          <button
            type="button"
            onClick={handleSave}
            disabled={isSaving || !value || isNaN(parseInt(value, 10)) || parseInt(value, 10) < 1}
            className={cn(
              'inline-flex items-center gap-1 rounded-md px-2.5 py-1.5 text-xs font-medium transition-colors',
              'bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50',
            )}
          >
            {isSaving ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <Save className="h-3.5 w-3.5" />}
            Set
          </button>
          {limit && (
            <button
              type="button"
              onClick={() => onRemove(status)}
              disabled={isRemoving}
              className="inline-flex items-center gap-1 rounded-md px-2.5 py-1.5 text-xs font-medium text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive disabled:opacity-50"
            >
              {isRemoving ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <Trash2 className="h-3.5 w-3.5" />
              )}
              Remove
            </button>
          )}
        </div>
      )}
    </div>
  );
}

export function WIPLimitsManager({ projectId, canManage }: Props) {
  const { data: wipLimits, isLoading } = useWIPLimits(projectId);
  const setLimit = useSetWIPLimit(projectId);
  const deleteLimit = useDeleteWIPLimit(projectId);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <Loader2 className="h-5 w-5 animate-spin" />
      </div>
    );
  }

  const limitMap = new Map((wipLimits ?? []).map((l) => [l.status, l]));

  return (
    <div className="space-y-4">
      <h3 className="flex items-center gap-2 text-sm font-semibold text-foreground">
        <Gauge className="h-4 w-4" />
        WIP Limits
      </h3>

      <div className="space-y-2">
        {statuses.map((status) => (
          <LimitRow
            key={status}
            status={status}
            limit={limitMap.get(status)}
            canManage={canManage}
            onSave={(s, max) => setLimit.mutate({ status: s, maxTasks: max })}
            onRemove={(s) => deleteLimit.mutate(s)}
            isSaving={setLimit.isPending}
            isRemoving={deleteLimit.isPending}
          />
        ))}
      </div>

      {!canManage && (
        <p className="text-xs text-muted-foreground">Only admins and managers can modify WIP limits.</p>
      )}
    </div>
  );
}
