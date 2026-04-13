import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { Calendar, GripVertical, ArrowRight, AlertOctagon } from 'lucide-react';
import type { Task, WorkspaceMember } from '@/types';
import { priorityLabels, priorityColors, isOverdue } from '@utils/format';
import { cn } from '@utils/cn';
import { format, parseISO } from 'date-fns';

function formatShortDate(date: string): string {
  return format(parseISO(date), 'MMM d');
}

interface Props {
  task: Task;
  members: WorkspaceMember[];
  onClick: () => void;
}

export function TaskCard({ task, members, onClick }: Props) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: task.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const assignee = members.find((m) => m.user_id === task.assignee_id);
  const overdue = isOverdue(task.due_date, task.status);

  const hasStartDate = !!task.start_date;
  const hasDueDate = !!task.due_date;
  const showDateRange = hasStartDate && hasDueDate;

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        'rounded-lg border border-border bg-card p-3 shadow-sm cursor-pointer',
        'hover:border-primary/30 transition-all',
        isDragging && 'opacity-50 shadow-lg rotate-2'
      )}
      onClick={onClick}
    >
      <div className="flex items-start gap-2">
        <button {...attributes} {...listeners} aria-label="Reorder task" className="mt-0.5 cursor-grab text-muted-foreground hover:text-foreground" onClick={(e) => e.stopPropagation()}>
          <GripVertical className="h-4 w-4" />
        </button>
        <div className="flex-1 min-w-0">
          {task.task_key && (
            <span className="text-xs font-mono text-muted-foreground">{task.task_key}</span>
          )}
          <p className="text-sm font-medium text-foreground truncate">{task.title}</p>
          <div className="mt-2 flex flex-wrap items-center gap-2">
            <span className={cn('rounded-full px-2 py-0.5 text-xs font-medium', priorityColors[task.priority])}>
              {priorityLabels[task.priority]}
            </span>
            {showDateRange ? (
              <span className={cn('flex items-center gap-1 text-xs', overdue ? 'text-destructive' : 'text-muted-foreground')}>
                <Calendar className="h-3 w-3" />
                {formatShortDate(task.start_date!)}
                <ArrowRight className="h-2.5 w-2.5" />
                {formatShortDate(task.due_date!)}
              </span>
            ) : hasDueDate ? (
              <span className={cn('flex items-center gap-1 text-xs', overdue ? 'text-destructive' : 'text-muted-foreground')}>
                <Calendar className="h-3 w-3" />
                {formatShortDate(task.due_date!)}
              </span>
            ) : hasStartDate ? (
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Calendar className="h-3 w-3" />
                {formatShortDate(task.start_date!)}
              </span>
            ) : null}
          </div>
          {task.status === 'blocked' && (task.blocked_by_task || task.blocked_reason) && (
            <div className="mt-2 flex items-center gap-1 text-xs text-red-600 dark:text-red-400">
              <AlertOctagon className="h-3 w-3" />
              <span className="truncate">
                {task.blocked_by_task ? `Blocked by ${task.blocked_by_task}` : 'Blocked'}
              </span>
            </div>
          )}
          {assignee && (
            <div className="mt-2 flex items-center gap-1.5">
              <div className="flex h-5 w-5 items-center justify-center rounded-full bg-primary/10 text-[10px] font-medium text-primary">
                {assignee.user_name.charAt(0).toUpperCase()}
              </div>
              <span className="text-xs text-muted-foreground">{assignee.user_name}</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
