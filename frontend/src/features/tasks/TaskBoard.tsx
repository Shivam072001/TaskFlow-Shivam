import { DndContext, DragOverlay, closestCorners, type DragEndEvent, type DragStartEvent } from '@dnd-kit/core';
import { useState, useCallback, useMemo } from 'react';
import type { Task, TaskStatus, WIPLimit, WorkspaceMember } from '@/types';
import { TaskColumn } from './TaskColumn';
import { TaskCard } from './TaskCard';

const columns: TaskStatus[] = ['todo', 'in_progress', 'blocked', 'done'];
const noop = () => {};

interface Props {
  tasks: Task[];
  members: WorkspaceMember[];
  onStatusChange: (taskId: string, newStatus: TaskStatus) => void;
  onTaskClick: (task: Task) => void;
  wipLimits?: Map<TaskStatus, WIPLimit>;
}

export function TaskBoard({ tasks, members, onStatusChange, onTaskClick, wipLimits }: Props) {
  const [activeId, setActiveId] = useState<string | null>(null);
  const activeTask = tasks.find((t) => t.id === activeId);

  const tasksByStatus = useMemo(() => {
    const grouped: Record<TaskStatus, Task[]> = { todo: [], in_progress: [], blocked: [], done: [] };
    for (const t of tasks) grouped[t.status].push(t);
    return grouped;
  }, [tasks]);

  const handleDragStart = useCallback((event: DragStartEvent) => {
    setActiveId(event.active.id as string);
  }, []);

  const handleDragEnd = useCallback(
    (event: DragEndEvent) => {
      setActiveId(null);
      const { active, over } = event;
      if (!over) return;

      const taskId = active.id as string;
      const task = tasks.find((t) => t.id === taskId);
      if (!task) return;

      let newStatus: TaskStatus | undefined;

      if (columns.includes(over.id as TaskStatus)) {
        newStatus = over.id as TaskStatus;
      } else {
        const overTask = tasks.find((t) => t.id === over.id);
        if (overTask) newStatus = overTask.status;
      }

      if (newStatus && newStatus !== task.status) {
        onStatusChange(taskId, newStatus);
      }
    },
    [tasks, onStatusChange],
  );

  return (
    <DndContext collisionDetection={closestCorners} onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        {columns.map((status) => (
          <TaskColumn
            key={status}
            status={status}
            tasks={tasksByStatus[status]}
            members={members}
            onTaskClick={onTaskClick}
            wipLimit={wipLimits?.get(status)?.max_tasks}
          />
        ))}
      </div>
      <DragOverlay>
        {activeTask ? <TaskCard task={activeTask} members={members} onClick={noop} /> : null}
      </DragOverlay>
    </DndContext>
  );
}
