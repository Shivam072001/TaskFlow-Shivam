import type React from 'react';
import { useState, useMemo } from 'react';
import { X, Loader2, ChevronDown, ChevronUp } from 'lucide-react';
import type { TaskStatus } from '@/types';
import { statusLabels } from '@utils/format';
import { cn } from '@utils/cn';

export interface WIPLimitInput {
  status: TaskStatus;
  maxTasks: number;
}

interface Props {
  open: boolean;
  existingPrefixes: string[];
  onClose: () => void;
  onSubmit: (data: { name: string; prefix: string; description: string; wipLimits: WIPLimitInput[] }) => void;
  isLoading: boolean;
}

const PREFIX_MAX = 6;
const PREFIX_MIN = 2;
const PREFIX_REGEX = /^[A-Z][A-Z0-9]*$/;
const wipStatuses: TaskStatus[] = ['todo', 'in_progress', 'blocked', 'done'];

export function CreateProjectDialog({ open, existingPrefixes, onClose, onSubmit, isLoading }: Props) {
  const [name, setName] = useState('');
  const [prefix, setPrefix] = useState('');
  const [description, setDescription] = useState('');
  const [showWip, setShowWip] = useState(false);
  const [wipValues, setWipValues] = useState<Record<string, string>>({});

  const takenPrefixes = useMemo(
    () => new Set(existingPrefixes.map((p) => p.toUpperCase())),
    [existingPrefixes],
  );

  if (!open) return null;

  const prefixUpper = prefix.toUpperCase();
  const prefixError = (() => {
    if (!prefix) return null;
    if (!PREFIX_REGEX.test(prefixUpper)) return 'Must start with a letter and contain only A-Z, 0-9';
    if (prefix.length < PREFIX_MIN) return `Minimum ${PREFIX_MIN} characters`;
    if (takenPrefixes.has(prefixUpper)) return `"${prefixUpper}" is already taken in this organization`;
    return null;
  })();
  const prefixValid = prefix.length >= PREFIX_MIN && !prefixError;

  const handleWipChange = (status: string, value: string) => {
    setWipValues((prev) => ({ ...prev, [status]: value }));
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!name.trim() || !prefixValid) return;

    const wipLimits: WIPLimitInput[] = [];
    for (const [status, val] of Object.entries(wipValues)) {
      const num = parseInt(val, 10);
      if (!isNaN(num) && num > 0) {
        wipLimits.push({ status: status as TaskStatus, maxTasks: num });
      }
    }

    onSubmit({
      name: name.trim(),
      prefix: prefixUpper,
      description: description.trim(),
      wipLimits,
    });
    setName('');
    setPrefix('');
    setDescription('');
    setWipValues({});
    setShowWip(false);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50" onClick={onClose}>
      <div
        className="w-full max-w-md max-h-[90vh] overflow-y-auto rounded-xl border border-border bg-card p-6 shadow-lg"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-foreground">New Project</h2>
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
              className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
              placeholder="Website Redesign"
              autoFocus
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Prefix</label>
            <input
              type="text"
              value={prefix}
              onChange={(e) => setPrefix(e.target.value.toUpperCase().replace(/[^A-Z0-9]/g, '').slice(0, PREFIX_MAX))}
              className={cn(
                'w-full rounded-lg border bg-background px-3 py-2 text-sm font-mono uppercase focus:outline-none focus:ring-2',
                prefixError
                  ? 'border-destructive focus:ring-destructive/50'
                  : 'border-border focus:ring-primary/50',
              )}
              placeholder="WR"
              maxLength={PREFIX_MAX}
            />
            {prefixError ? (
              <p className="mt-1 text-xs text-destructive">{prefixError}</p>
            ) : (
              <p className="mt-1 text-xs text-muted-foreground">
                Used for task keys (e.g. {prefixUpper || 'WR'}-1, {prefixUpper || 'WR'}-2). {PREFIX_MIN}-{PREFIX_MAX} uppercase letters/numbers.
              </p>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1.5">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-primary/50"
              placeholder="Optional description..."
            />
          </div>

          {/* WIP Limits (optional) */}
          <div className="border-t border-border pt-4">
            <button
              type="button"
              onClick={() => setShowWip((o) => !o)}
              className="flex w-full items-center justify-between text-sm font-medium text-foreground"
            >
              <span>WIP Limits <span className="font-normal text-muted-foreground">(optional)</span></span>
              {showWip ? <ChevronUp className="h-4 w-4 text-muted-foreground" /> : <ChevronDown className="h-4 w-4 text-muted-foreground" />}
            </button>
            {showWip && (
              <div className="mt-3 space-y-2">
                <p className="text-xs text-muted-foreground">
                  Set max tasks allowed per column. Leave blank for no limit.
                </p>
                {wipStatuses.map((status) => (
                  <div key={status} className="flex items-center gap-3">
                    <span className="w-28 text-sm text-foreground">{statusLabels[status]}</span>
                    <input
                      type="number"
                      min={1}
                      value={wipValues[status] ?? ''}
                      onChange={(e) => handleWipChange(status, e.target.value)}
                      placeholder="No limit"
                      className="w-24 rounded-md border border-border bg-background px-2 py-1.5 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                    />
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="flex justify-end gap-2">
            <button type="button" onClick={onClose} className="rounded-lg px-4 py-2 text-sm text-muted-foreground hover:bg-accent">
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading || !name.trim() || !prefixValid}
              className={cn('rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50')}
            >
              {isLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
