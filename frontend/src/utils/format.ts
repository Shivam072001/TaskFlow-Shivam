import { format, formatDistanceToNow, isPast, parseISO } from 'date-fns';
import type { SelectOption } from '@components/ui/Select';
import type { TaskPriority, TaskStatus, WorkspaceRole } from '@/types';

export function formatDate(date: string | null): string {
  if (!date) return '';
  return format(parseISO(date), 'MMM d, yyyy');
}

export function formatRelative(date: string): string {
  return formatDistanceToNow(parseISO(date), { addSuffix: true });
}

export function isOverdue(date: string | null, status: TaskStatus): boolean {
  if (!date || status === 'done') return false;
  return isPast(parseISO(date));
}

export const statusLabels: Record<TaskStatus, string> = {
  todo: 'To Do',
  in_progress: 'In Progress',
  done: 'Done',
  blocked: 'Blocked',
};

export const priorityLabels: Record<TaskPriority, string> = {
  low: 'Low',
  medium: 'Medium',
  high: 'High',
};

export const priorityColors: Record<TaskPriority, string> = {
  low: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
  medium: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
  high: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
};

export const statusOptions: SelectOption[] = [
  { value: 'todo', label: 'To Do' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'blocked', label: 'Blocked' },
  { value: 'done', label: 'Done' },
];

export const priorityOptions: SelectOption[] = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
];

export const roleLabels: Record<WorkspaceRole, string> = {
  owner: 'Owner',
  admin: 'Admin',
  manager: 'Manager',
  lead: 'Lead',
  member: 'Member',
  guest: 'Guest',
  viewer: 'Viewer',
};

export const rolePower: Record<WorkspaceRole, number> = {
  owner: 100,
  admin: 80,
  manager: 60,
  lead: 40,
  member: 20,
  guest: 10,
  viewer: 5,
};
