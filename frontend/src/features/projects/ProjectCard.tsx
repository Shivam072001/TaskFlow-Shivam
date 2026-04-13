import { Link } from 'react-router-dom';
import { Kanban } from 'lucide-react';
import type { Project } from '@/types';
import { formatRelative } from '@utils/format';
import { cn } from '@utils/cn';

interface Props {
  project: Project;
  workspaceId: string;
  orgSlug: string;
}

export function ProjectCard({ project, workspaceId, orgSlug }: Props) {
  return (
    <Link
      to={`/org/${orgSlug}/workspaces/${workspaceId}/projects/${project.id}`}
      className={cn(
        'block rounded-xl border border-border bg-card p-4',
        'hover:border-primary/50 hover:shadow-sm transition-all'
      )}
    >
      <div className="flex items-center gap-2 mb-2">
        <Kanban className="h-4 w-4 text-primary" />
        <h4 className="font-medium text-foreground">{project.name}</h4>
        {project.prefix && (
          <span className="rounded bg-muted px-1.5 py-0.5 text-xs font-mono font-medium text-muted-foreground">
            {project.prefix}
          </span>
        )}
      </div>
      {project.description && (
        <p className="text-sm text-muted-foreground line-clamp-2 mb-2">{project.description}</p>
      )}
      <p className="text-xs text-muted-foreground">Created {formatRelative(project.created_at)}</p>
    </Link>
  );
}
