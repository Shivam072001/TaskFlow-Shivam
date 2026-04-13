import type React from 'react';
import { useState } from 'react';
import { X, Loader2 } from 'lucide-react';
import { cn } from '@utils/cn';

interface Props {
  open: boolean;
  onClose: () => void;
  onSubmit: (name: string, description: string) => void;
  isLoading: boolean;
}

export function CreateWorkspaceDialog({ open, onClose, onSubmit, isLoading }: Props) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  if (!open) return null;

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (name.trim()) onSubmit(name.trim(), description.trim());
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <div
        className="w-full max-w-md rounded-xl border border-border bg-card p-6 shadow-lg"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-foreground">Create Workspace</h2>
          <button onClick={onClose} className="text-muted-foreground hover:text-foreground">
            <X className="h-4 w-4" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className={cn(
                'w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground',
                'focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary'
              )}
              placeholder="Engineering Team"
              autoFocus
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className={cn(
                'w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground resize-none',
                'focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary'
              )}
              placeholder="Optional description..."
            />
          </div>
          <div className="flex justify-end gap-2">
            <button type="button" onClick={onClose} className="rounded-lg px-4 py-2 text-sm text-muted-foreground hover:bg-accent">
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading || !name.trim()}
              className={cn(
                'rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground',
                'hover:bg-primary/90 disabled:opacity-50'
              )}
            >
              {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
