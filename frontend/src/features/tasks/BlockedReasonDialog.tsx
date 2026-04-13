import type React from 'react';
import { useState } from 'react';
import { AlertOctagon, X } from 'lucide-react';
import { cn } from '@utils/cn';

interface Props {
  open: boolean;
  onClose: () => void;
  onConfirm: (blockedByTask: string, blockedReason: string) => void;
}

export function BlockedReasonDialog({ open, onClose, onConfirm }: Props) {
  const [blockedByTask, setBlockedByTask] = useState('');
  const [blockedReason, setBlockedReason] = useState('');

  if (!open) return null;

  const isValid = blockedByTask.trim() !== '' || blockedReason.trim() !== '';

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!isValid) return;
    onConfirm(blockedByTask.trim(), blockedReason.trim());
    setBlockedByTask('');
    setBlockedReason('');
  };

  const handleClose = () => {
    setBlockedByTask('');
    setBlockedReason('');
    onClose();
  };

  return (
    <div className="fixed inset-0 z-[110] flex items-center justify-center bg-black/50 backdrop-blur-[2px]" onClick={handleClose}>
      <div
        className="w-full max-w-sm rounded-xl border border-border bg-card p-6 shadow-xl animate-in zoom-in-95 duration-150"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-start gap-4">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400">
            <AlertOctagon className="h-5 w-5" />
          </div>
          <div className="flex-1">
            <div className="flex items-start justify-between">
              <h3 className="text-base font-semibold text-foreground">Mark as Blocked</h3>
              <button onClick={handleClose} className="rounded-md p-0.5 text-muted-foreground hover:text-foreground">
                <X className="h-4 w-4" />
              </button>
            </div>
            <p className="mt-1.5 text-sm text-muted-foreground">
              Provide at least one: the blocking task ID or a reason.
            </p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="mt-5 space-y-4">
          <div>
            <label className="mb-1.5 block text-sm font-medium text-foreground">Blocking Task ID</label>
            <input
              type="text"
              value={blockedByTask}
              onChange={(e) => setBlockedByTask(e.target.value)}
              placeholder="e.g. IP-42"
              className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm font-mono text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
          </div>
          <div>
            <label className="mb-1.5 block text-sm font-medium text-foreground">Reason</label>
            <textarea
              value={blockedReason}
              onChange={(e) => setBlockedReason(e.target.value)}
              rows={3}
              placeholder="Why is this task blocked?"
              className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground resize-none placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/50"
            />
          </div>

          <div className="flex justify-end gap-2">
            <button type="button" onClick={handleClose}
              className="rounded-lg px-4 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted">
              Cancel
            </button>
            <button type="submit" disabled={!isValid}
              className={cn(
                'rounded-lg px-4 py-2 text-sm font-medium text-white transition-colors',
                'bg-red-600 hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed',
              )}>
              Mark Blocked
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
