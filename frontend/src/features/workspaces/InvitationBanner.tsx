import { Check, X, Mail, Loader2 } from 'lucide-react';
import type { WorkspaceInvitation } from '@/types';
import { roleLabels } from '@utils/format';
import { Badge } from '@components/ui/Badge';
import { cn } from '@utils/cn';

interface Props {
  invitations: WorkspaceInvitation[];
  onRespond: (invitationId: string, accept: boolean) => void;
  isResponding: boolean;
}

export function InvitationBanner({ invitations, onRespond, isResponding }: Props) {
  const pending = invitations.filter((inv) => inv.status === 'pending');

  if (pending.length === 0) return null;

  return (
    <div className="space-y-3">
      <h3 className="flex items-center gap-2 text-sm font-semibold text-foreground">
        <Mail className="h-4 w-4" />
        Pending Invitations
        <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs font-normal text-primary">
          {pending.length}
        </span>
      </h3>

      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {pending.map((inv) => (
          <div
            key={inv.id}
            className="rounded-lg border border-border bg-card p-4 shadow-sm transition-colors hover:border-primary/30"
          >
            <div className="mb-3">
              <p className="text-sm font-medium text-foreground">{inv.workspace_name}</p>
              <p className="mt-0.5 text-xs text-muted-foreground">
                Invited by {inv.inviter_name}
              </p>
            </div>

            <div className="mb-3">
              <Badge variant={inv.role}>{roleLabels[inv.role]}</Badge>
            </div>

            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => onRespond(inv.id, true)}
                disabled={isResponding}
                className={cn(
                  'inline-flex flex-1 items-center justify-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors',
                  'bg-green-600 text-white hover:bg-green-700 dark:bg-green-700 dark:hover:bg-green-600',
                  'disabled:opacity-50',
                )}
              >
                {isResponding ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <Check className="h-3.5 w-3.5" />}
                Accept
              </button>
              <button
                type="button"
                onClick={() => onRespond(inv.id, false)}
                disabled={isResponding}
                className={cn(
                  'inline-flex flex-1 items-center justify-center gap-1.5 rounded-md px-3 py-1.5 text-xs font-medium transition-colors',
                  'bg-red-600 text-white hover:bg-red-700 dark:bg-red-700 dark:hover:bg-red-600',
                  'disabled:opacity-50',
                )}
              >
                {isResponding ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <X className="h-3.5 w-3.5" />}
                Decline
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
