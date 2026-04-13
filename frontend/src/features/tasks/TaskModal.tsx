import type React from 'react';
import { useState, type ReactNode } from 'react';
import { X, Loader2, Trash2, AlertOctagon, Pencil, Calendar, ArrowRight } from 'lucide-react';
import type { Task, WorkspaceMember, CustomFieldDefinition, TaskStatus } from '@/types';
import { RichTextEditor } from '@components/ui/RichTextEditor';
import { RichTextDisplay } from '@components/ui/RichTextDisplay';
import { Select } from '@components/ui/Select';
import { priorityLabels, priorityColors, formatDate, statusOptions, priorityOptions } from '@utils/format';
import { cn } from '@utils/cn';

interface Props {
  open: boolean;
  task: Task | null;
  members: WorkspaceMember[];
  customFields: CustomFieldDefinition[];
  customFieldValues: Record<string, string>;
  onClose: () => void;
  onSave: (data: {
    title: string;
    description: string;
    status?: string;
    priority: string;
    assignee_id: string;
    start_date: string;
    due_date: string;
    custom_fields?: Record<string, string>;
  }) => void;
  onStatusChange?: (taskId: string, status: TaskStatus) => void;
  onDelete?: () => void;
  isSaving: boolean;
  children?: ReactNode;
}

export function TaskModal({
  open,
  task,
  members,
  customFields,
  customFieldValues,
  onClose,
  onSave,
  onStatusChange,
  onDelete,
  isSaving,
  children,
}: Props) {
  const isEdit = !!task;
  const [mode, setMode] = useState<'view' | 'edit'>(isEdit ? 'view' : 'edit');

  const [title, setTitle] = useState(task?.title ?? '');
  const [description, setDescription] = useState(task?.description ?? '');
  const [status, setStatus] = useState<string>(task?.status ?? 'todo');
  const [priority, setPriority] = useState<string>(task?.priority ?? 'medium');
  const [assigneeId, setAssigneeId] = useState(task?.assignee_id ?? '');
  const [startDate, setStartDate] = useState(task?.start_date ?? '');
  const [dueDate, setDueDate] = useState(task?.due_date ?? '');
  const [cfValues, setCfValues] = useState<Record<string, string>>(customFieldValues);

  if (!open) return null;

  const assignee = members.find((m) => m.user_id === (mode === 'view' ? task?.assignee_id : assigneeId));

  const switchToEdit = () => {
    setTitle(task?.title ?? '');
    setDescription(task?.description ?? '');
    setStatus(task?.status ?? 'todo');
    setPriority(task?.priority ?? 'medium');
    setAssigneeId(task?.assignee_id ?? '');
    setStartDate(task?.start_date ?? '');
    setDueDate(task?.due_date ?? '');
    setCfValues(customFieldValues);
    setMode('edit');
  };

  const handleCfChange = (fieldId: string, value: string) => {
    setCfValues((prev) => ({ ...prev, [fieldId]: value }));
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!title.trim()) return;
    onSave({
      title: title.trim(),
      description,
      ...(task ? { status } : {}),
      priority,
      assignee_id: assigneeId,
      start_date: startDate,
      due_date: dueDate,
      custom_fields: cfValues,
    });
  };

  const requiredFieldsMissing = customFields.some(
    (f) => f.required && !(cfValues[f.id] ?? '').trim(),
  );
  const canSubmit = !!title.trim() && !isSaving && !requiredFieldsMissing;

  const inputClasses =
    'w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-primary/50';

  const assigneeOptions = [
    { value: '', label: 'Unassigned' },
    ...members.map((m) => ({ value: m.user_id, label: m.user_name })),
  ];

  // ── View Mode ──
  if (mode === 'view' && task) {
    return (
      <div className="fixed inset-0 z-50 flex justify-end bg-black/50" onClick={onClose}>
        <div
          className="h-full w-full max-w-2xl overflow-y-auto border-l border-border bg-card shadow-xl"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Header */}
          <div className="sticky top-0 z-10 flex items-center justify-between border-b border-border bg-card px-6 py-4">
            <div className="min-w-0 flex-1">
              {task.task_key && (
                <span className="text-xs font-mono text-muted-foreground">{task.task_key}</span>
              )}
              <h2 className="text-lg font-semibold text-foreground truncate">{task.title}</h2>
            </div>
            <div className="flex items-center gap-2 ml-4">
              <button
                onClick={switchToEdit}
                className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
              >
                <Pencil className="h-3.5 w-3.5" /> Edit
              </button>
              <button onClick={onClose} className="rounded-md p-1.5 text-muted-foreground hover:text-foreground">
                <X className="h-5 w-5" />
              </button>
            </div>
          </div>

          <div className="px-6 py-5 space-y-6">
            {/* Status + Quick Change */}
            <div className="flex flex-wrap items-center gap-3">
              <Select
                value={task.status}
                onChange={(val) => onStatusChange?.(task.id, val as TaskStatus)}
                options={statusOptions}
                className="w-40"
              />
              <span className={cn('rounded-full px-2.5 py-0.5 text-xs font-medium', priorityColors[task.priority])}>
                {priorityLabels[task.priority]}
              </span>
            </div>

            {/* Blocked Banner */}
            {task.status === 'blocked' && (task.blocked_by_task || task.blocked_reason) && (
              <div className="rounded-lg border border-red-200 bg-red-50 p-3 dark:border-red-900/50 dark:bg-red-950/30">
                <div className="flex items-center gap-2 text-sm font-medium text-red-700 dark:text-red-400">
                  <AlertOctagon className="h-4 w-4" />
                  Blocked
                </div>
                {task.blocked_by_task && (
                  <p className="mt-1.5 text-sm text-red-600 dark:text-red-400">
                    <span className="font-medium">Blocked by:</span>{' '}
                    <span className="font-mono">{task.blocked_by_task}</span>
                  </p>
                )}
                {task.blocked_reason && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                    <span className="font-medium">Reason:</span> {task.blocked_reason}
                  </p>
                )}
              </div>
            )}

            {/* Meta Row */}
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div className="space-y-1">
                <span className="text-xs font-medium uppercase tracking-wider text-muted-foreground">Assignee</span>
                <div className="flex items-center gap-2 text-foreground">
                  {assignee ? (
                    <>
                      <div className="flex h-6 w-6 items-center justify-center rounded-full bg-primary/10 text-xs font-medium text-primary">
                        {assignee.user_name.charAt(0).toUpperCase()}
                      </div>
                      {assignee.user_name}
                    </>
                  ) : (
                    <span className="text-muted-foreground">Unassigned</span>
                  )}
                </div>
              </div>
              <div className="space-y-1">
                <span className="text-xs font-medium uppercase tracking-wider text-muted-foreground">Dates</span>
                <div className="flex items-center gap-1.5 text-foreground">
                  <Calendar className="h-3.5 w-3.5 text-muted-foreground" />
                  {task.start_date && task.due_date ? (
                    <>
                      {formatDate(task.start_date)}
                      <ArrowRight className="h-3 w-3 text-muted-foreground" />
                      {formatDate(task.due_date)}
                    </>
                  ) : task.due_date ? (
                    <>Due {formatDate(task.due_date)}</>
                  ) : task.start_date ? (
                    <>Starts {formatDate(task.start_date)}</>
                  ) : (
                    <span className="text-muted-foreground">No dates set</span>
                  )}
                </div>
              </div>
            </div>

            {/* Description */}
            <div>
              <h3 className="mb-2 text-xs font-medium uppercase tracking-wider text-muted-foreground">Description</h3>
              {task.description ? (
                <div className="rounded-lg border border-border bg-muted/30 px-4 py-3">
                  <RichTextDisplay content={task.description} />
                </div>
              ) : (
                <p className="text-sm italic text-muted-foreground">No description provided.</p>
              )}
            </div>

            {/* Custom Field Values */}
            {customFields.length > 0 && (
              <div>
                <h3 className="mb-3 text-xs font-medium uppercase tracking-wider text-muted-foreground">
                  Custom Fields
                </h3>
                <div className="grid grid-cols-2 gap-3">
                  {customFields.map((field) => (
                    <div key={field.id} className="rounded-md border border-border px-3 py-2">
                      <span className="text-xs text-muted-foreground">{field.name}</span>
                      <p className="text-sm font-medium text-foreground">
                        {customFieldValues[field.id] || <span className="text-muted-foreground italic">—</span>}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Discussion (children) */}
            {children}
          </div>
        </div>
      </div>
    );
  }

  // ── Edit / Create Mode ──
  return (
    <div className="fixed inset-0 z-50 flex justify-end bg-black/50" onClick={onClose}>
      <div
        className="h-full w-full max-w-lg overflow-y-auto border-l border-border bg-card p-6 shadow-xl"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="mb-6 flex items-center justify-between">
          <div>
            {isEdit && task?.task_key && (
              <span className="text-xs font-mono text-muted-foreground">{task.task_key}</span>
            )}
            <h2 className="text-lg font-semibold text-foreground">
              {isEdit ? 'Edit Task' : 'New Task'}
            </h2>
          </div>
          <div className="flex items-center gap-2">
            {isEdit && (
              <button
                onClick={() => setMode('view')}
                className="rounded-lg px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent"
              >
                Back
              </button>
            )}
            <button onClick={onClose} className="text-muted-foreground hover:text-foreground">
              <X className="h-5 w-5" />
            </button>
          </div>
        </div>

        {isEdit && task?.status === 'blocked' && (task.blocked_by_task || task.blocked_reason) && (
          <div className="mb-5 rounded-lg border border-red-200 bg-red-50 p-3 dark:border-red-900/50 dark:bg-red-950/30">
            <div className="flex items-center gap-2 text-sm font-medium text-red-700 dark:text-red-400">
              <AlertOctagon className="h-4 w-4" />
              Blocked
            </div>
            {task.blocked_by_task && (
              <p className="mt-1.5 text-sm text-red-600 dark:text-red-400">
                <span className="font-medium">Blocked by:</span>{' '}
                <span className="font-mono">{task.blocked_by_task}</span>
              </p>
            )}
            {task.blocked_reason && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                <span className="font-medium">Reason:</span> {task.blocked_reason}
              </p>
            )}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="mb-1.5 block text-sm font-medium text-foreground">Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className={inputClasses}
              placeholder="Task title"
              autoFocus
            />
          </div>

          <div>
            <label className="mb-1.5 block text-sm font-medium text-foreground">Description</label>
            <RichTextEditor
              value={description}
              onChange={setDescription}
              placeholder="Add a description..."
            />
          </div>

          {isEdit && (
            <div>
              <label className="mb-1.5 block text-sm font-medium text-foreground">Status</label>
              <Select value={status} onChange={setStatus} options={statusOptions} />
            </div>
          )}

          <div>
            <label className="mb-1.5 block text-sm font-medium text-foreground">Priority</label>
            <Select value={priority} onChange={setPriority} options={priorityOptions} />
          </div>

          <div>
            <label className="mb-1.5 block text-sm font-medium text-foreground">Assignee</label>
            <Select
              value={assigneeId}
              onChange={setAssigneeId}
              options={assigneeOptions}
              placeholder="Unassigned"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="mb-1.5 block text-sm font-medium text-foreground">Start Date</label>
              <input
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className={inputClasses}
              />
            </div>
            <div>
              <label className="mb-1.5 block text-sm font-medium text-foreground">Due Date</label>
              <input
                type="date"
                value={dueDate}
                onChange={(e) => setDueDate(e.target.value)}
                className={inputClasses}
              />
            </div>
          </div>

          {customFields.length > 0 && (
            <div className="space-y-4 border-t border-border pt-4">
              <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                Custom Fields
              </p>
              {customFields.map((field) => (
                <div key={field.id}>
                  <label className="mb-1.5 block text-sm font-medium text-foreground">
                    {field.name}
                    {field.required && <span className="ml-0.5 text-destructive">*</span>}
                  </label>

                  {field.field_type === 'text' && (
                    <input
                      type="text"
                      value={cfValues[field.id] ?? ''}
                      onChange={(e) => handleCfChange(field.id, e.target.value)}
                      className={inputClasses}
                    />
                  )}

                  {field.field_type === 'number' && (
                    <input
                      type="number"
                      value={cfValues[field.id] ?? ''}
                      onChange={(e) => handleCfChange(field.id, e.target.value)}
                      className={inputClasses}
                    />
                  )}

                  {field.field_type === 'select' && (
                    <Select
                      value={cfValues[field.id] ?? ''}
                      onChange={(v) => handleCfChange(field.id, v)}
                      options={field.options.map((o) => ({ value: o, label: o }))}
                      placeholder={`Select ${field.name.toLowerCase()}…`}
                    />
                  )}
                </div>
              ))}
            </div>
          )}

          <div className="flex items-center gap-2 pt-2">
            <button
              type="submit"
              disabled={!canSubmit}
              className="rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
            >
              {isSaving ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : isEdit ? (
                'Update'
              ) : (
                'Create'
              )}
            </button>
            <button
              type="button"
              onClick={isEdit ? () => setMode('view') : onClose}
              className="rounded-lg px-4 py-2.5 text-sm text-muted-foreground hover:bg-accent"
            >
              Cancel
            </button>
            {isEdit && onDelete && (
              <button
                type="button"
                onClick={onDelete}
                className="ml-auto rounded-lg px-4 py-2.5 text-sm text-destructive hover:bg-destructive/10"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            )}
          </div>
        </form>

        {children}
      </div>
    </div>
  );
}
