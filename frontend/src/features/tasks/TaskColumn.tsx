import { useDroppable } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import type { Task, TaskStatus, WorkspaceMember } from '@/types';
import { statusLabels } from '@utils/format';
import { cn } from '@utils/cn';
import { TaskCard } from './TaskCard';

const columnColors: Record<TaskStatus, string> = {
  todo: 'bg-yellow-500',
  in_progress: 'bg-blue-500',
  blocked: 'bg-red-500',
  done: 'bg-green-500',
};

interface Props {
  status: TaskStatus;
  tasks: Task[];
  members: WorkspaceMember[];
  onTaskClick: (task: Task) => void;
  wipLimit?: number;
}

export function TaskColumn({ status, tasks, members, onTaskClick, wipLimit }: Props) {
  const { setNodeRef, isOver } = useDroppable({ id: status });
  const atCapacity = wipLimit !== undefined && tasks.length >= wipLimit;

  return (
    <div
      className={cn(
        'flex flex-col rounded-xl bg-muted/50 p-3 min-h-[200px]',
        isOver && 'ring-2 ring-primary/30',
        atCapacity && 'ring-1 ring-destructive/40 bg-destructive/5',
      )}
    >
      <div className="flex items-center gap-2 mb-3 px-1">
        <div className={cn('h-2 w-2 rounded-full', columnColors[status])} />
        <h3 className="text-sm font-semibold text-foreground">{statusLabels[status]}</h3>
        <span
          className={cn(
            'ml-auto rounded-full px-2 py-0.5 text-xs font-medium',
            atCapacity
              ? 'bg-destructive/10 text-destructive'
              : 'bg-muted text-muted-foreground',
          )}
        >
          {tasks.length}{wipLimit !== undefined ? `/${wipLimit}` : ''}
        </span>
      </div>
      <div ref={setNodeRef} className="flex flex-col gap-2 flex-1">
        <SortableContext items={tasks.map((t) => t.id)} strategy={verticalListSortingStrategy}>
          {tasks.map((task) => (
            <TaskCard key={task.id} task={task} members={members} onClick={() => onTaskClick(task)} />
          ))}
        </SortableContext>
        {tasks.length === 0 && (
          <div className="flex flex-1 items-center justify-center rounded-lg border border-dashed border-border p-4">
            <p className="text-xs text-muted-foreground">No tasks</p>
          </div>
        )}
      </div>
    </div>
  );
}
