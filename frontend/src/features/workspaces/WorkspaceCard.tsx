import { Link } from 'react-router-dom';
import { FolderOpen, Users } from 'lucide-react';
import type { Workspace } from '@/types';
import { formatRelative } from '@utils/format';
import { cn } from '@utils/cn';

export function WorkspaceCard({ workspace, orgSlug }: { workspace: Workspace; orgSlug: string }) {
  return (
    <Link
      to={`/org/${orgSlug}/workspaces/${workspace.id}`}
      className={cn(
        'block rounded-xl border border-border bg-card p-5',
        'hover:border-primary/50 hover:shadow-md transition-all'
      )}
    >
      <div className="flex items-start justify-between">
        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
          <FolderOpen className="h-5 w-5 text-primary" />
        </div>
      </div>
      <h3 className="mt-3 font-semibold text-foreground">{workspace.name}</h3>
      {workspace.description && (
        <p className="mt-1 text-sm text-muted-foreground line-clamp-2">{workspace.description}</p>
      )}
      <div className="mt-3 flex items-center gap-3 text-xs text-muted-foreground">
        <span className="flex items-center gap-1">
          <Users className="h-3 w-3" />
          Members
        </span>
        <span>Created {formatRelative(workspace.created_at)}</span>
      </div>
    </Link>
  );
}
