import { useState, useMemo } from 'react';
import {
  Loader2, Search, X, CheckCircle2, Clock, AlertTriangle,
  CircleDot, Filter, CalendarDays, ListTodo,
} from 'lucide-react';
import { PageShell } from '@components/layout/PageShell';
import { Pagination } from '@components/ui/Pagination';
import { useMyTasks, useMyStats, useMyProjects } from '@hooks/useDashboard';
import { useAuth } from '@hooks/useAuth';
import { cn } from '@utils/cn';
import type { Task, TaskStatus, TaskPriority } from '@/types';

const STATUS_CONFIG: Record<TaskStatus, { label: string; color: string; icon: typeof Clock }> = {
  todo: { label: 'To Do', color: 'bg-yellow-500', icon: CircleDot },
  in_progress: { label: 'In Progress', color: 'bg-blue-500', icon: Clock },
  blocked: { label: 'Blocked', color: 'bg-red-500', icon: AlertTriangle },
  done: { label: 'Done', color: 'bg-green-500', icon: CheckCircle2 },
};

const PRIORITY_COLORS: Record<TaskPriority, string> = {
  high: 'text-red-500 bg-red-500/10',
  medium: 'text-amber-500 bg-amber-500/10',
  low: 'text-emerald-500 bg-emerald-500/10',
};

export function MyDashboardPage() {
  const { user } = useAuth();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [statusFilter, setStatusFilter] = useState('');
  const [priorityFilter, setPriorityFilter] = useState('');
  const [projectFilter, setProjectFilter] = useState('');
  const [search, setSearch] = useState('');
  const [dueAfter, setDueAfter] = useState('');
  const [dueBefore, setDueBefore] = useState('');
  const [showFilters, setShowFilters] = useState(false);

  const filters = useMemo(() => ({
    status: statusFilter || undefined,
    priority: priorityFilter || undefined,
    search: search || undefined,
    due_before: dueBefore || undefined,
    due_after: dueAfter || undefined,
    project_id: projectFilter || undefined,
    page,
    limit: pageSize,
  }), [statusFilter, priorityFilter, search, dueBefore, dueAfter, projectFilter, page, pageSize]);

  const { data: tasksData, isLoading: tasksLoading, isError: tasksError } = useMyTasks(filters);
  const { data: stats, isLoading: statsLoading, isError: statsError } = useMyStats();
  const { data: projects } = useMyProjects();

  const tasks = tasksData?.data ?? [];
  const total = tasksData?.meta?.total ?? 0;

  const activeFilterCount = [statusFilter, priorityFilter, projectFilter, dueAfter, dueBefore].filter(Boolean).length;

  const clearFilters = () => {
    setStatusFilter('');
    setPriorityFilter('');
    setProjectFilter('');
    setDueAfter('');
    setDueBefore('');
    setSearch('');
    setPage(1);
  };

  const formatDate = (d: string | null) => {
    if (!d) return '—';
    return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  };

  const isOverdue = (task: Task) =>
    task.due_date && task.status !== 'done' && new Date(task.due_date) < new Date(new Date().toDateString());

  const projectMap = useMemo(() => {
    const m = new Map<string, string>();
    projects?.forEach(p => m.set(p.id, p.name));
    return m;
  }, [projects]);

  return (
    <PageShell>
      <div className="space-y-6">
        {/* Header */}
        <div>
          <h1 className="text-2xl font-bold text-foreground">
            Welcome back, {user?.name?.split(' ')[0] ?? 'there'}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Here's an overview of all tasks assigned to you across every project.
          </p>
        </div>

        {/* Stat Cards */}
        {statsLoading ? (
          <div className="flex h-24 items-center justify-center">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : statsError ? (
          <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-4 text-center text-sm text-destructive">
            Failed to load stats. Please try again.
          </div>
        ) : stats ? (
          <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6">
            <StatCard label="Total" value={stats.total} icon={ListTodo} color="text-foreground" />
            <StatCard label="In Progress" value={stats.by_status.in_progress ?? 0} icon={Clock} color="text-blue-500" />
            <StatCard label="To Do" value={stats.by_status.todo ?? 0} icon={CircleDot} color="text-slate-400" />
            <StatCard label="Completed" value={stats.completed} icon={CheckCircle2} color="text-emerald-500" />
            <StatCard label="Overdue" value={stats.overdue} icon={AlertTriangle} color="text-red-500" />
            <StatCard label="Due This Week" value={stats.due_this_week} icon={CalendarDays} color="text-amber-500" />
          </div>
        ) : null}

        {/* Priority breakdown */}
        {stats && stats.total > 0 && (
          <div className="flex flex-wrap items-center gap-4 rounded-lg border border-border bg-card p-4">
            <span className="text-sm font-medium text-muted-foreground">By Priority:</span>
            {(['high', 'medium', 'low'] as TaskPriority[]).map(p => (
              <button
                key={p}
                onClick={() => { setPriorityFilter(priorityFilter === p ? '' : p); setPage(1); }}
                className={cn(
                  'inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium transition-all',
                  PRIORITY_COLORS[p],
                  priorityFilter === p && 'ring-2 ring-offset-1 ring-offset-background ring-current',
                )}
              >
                {p.charAt(0).toUpperCase() + p.slice(1)}: {stats.by_priority[p] ?? 0}
              </button>
            ))}
          </div>
        )}

        {/* Search + Filter Bar */}
        <div className="space-y-3">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type="text"
                placeholder="Search tasks by title or key..."
                value={search}
                onChange={e => { setSearch(e.target.value); setPage(1); }}
                className="w-full rounded-lg border border-border bg-card py-2 pl-10 pr-4 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              />
              {search && (
                <button onClick={() => setSearch('')} className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground">
                  <X className="h-4 w-4" />
                </button>
              )}
            </div>
            <button
              onClick={() => setShowFilters(!showFilters)}
              className={cn(
                'inline-flex items-center gap-2 rounded-lg border border-border px-4 py-2 text-sm font-medium transition-colors',
                showFilters || activeFilterCount > 0 ? 'bg-primary/10 text-primary border-primary/30' : 'bg-card text-muted-foreground hover:text-foreground',
              )}
            >
              <Filter className="h-4 w-4" />
              Filters
              {activeFilterCount > 0 && (
                <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-primary-foreground">
                  {activeFilterCount}
                </span>
              )}
            </button>
            {activeFilterCount > 0 && (
              <button onClick={clearFilters} className="text-sm text-muted-foreground hover:text-foreground">
                Clear all
              </button>
            )}
          </div>

          {/* Expanded Filters */}
          {showFilters && (
            <div className="grid grid-cols-1 gap-3 rounded-lg border border-border bg-card p-4 sm:grid-cols-2 lg:grid-cols-4">
              <FilterSelect label="Status" value={statusFilter} onChange={v => { setStatusFilter(v); setPage(1); }}
                options={[{ value: '', label: 'All statuses' }, { value: 'todo', label: 'To Do' }, { value: 'in_progress', label: 'In Progress' }, { value: 'blocked', label: 'Blocked' }, { value: 'done', label: 'Done' }]} />
              <FilterSelect label="Priority" value={priorityFilter} onChange={v => { setPriorityFilter(v); setPage(1); }}
                options={[{ value: '', label: 'All priorities' }, { value: 'high', label: 'High' }, { value: 'medium', label: 'Medium' }, { value: 'low', label: 'Low' }]} />
              <FilterSelect label="Project" value={projectFilter} onChange={v => { setProjectFilter(v); setPage(1); }}
                options={[{ value: '', label: 'All projects' }, ...(projects ?? []).map(p => ({ value: p.id, label: p.name }))]} />
              <div className="space-y-1">
                <label className="text-xs font-medium text-muted-foreground">Due Date Range</label>
                <div className="flex items-center gap-2">
                  <input type="date" value={dueAfter} onChange={e => { setDueAfter(e.target.value); setPage(1); }}
                    className="w-full rounded-md border border-border bg-background px-2 py-1.5 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring" />
                  <span className="text-xs text-muted-foreground">to</span>
                  <input type="date" value={dueBefore} onChange={e => { setDueBefore(e.target.value); setPage(1); }}
                    className="w-full rounded-md border border-border bg-background px-2 py-1.5 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring" />
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Task Table */}
        {tasksLoading ? (
          <div className="flex h-40 items-center justify-center">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : tasksError ? (
          <div className="rounded-lg border border-destructive/30 bg-destructive/5 p-6 text-center text-sm text-destructive">
            Failed to load tasks. Please try again.
          </div>
        ) : tasks.length === 0 ? (
          <div className="rounded-lg border border-border bg-card p-12 text-center">
            <ListTodo className="mx-auto h-10 w-10 text-muted-foreground/40" />
            <p className="mt-3 text-sm font-medium text-muted-foreground">
              {activeFilterCount > 0 || search ? 'No tasks match your filters.' : 'No tasks assigned to you yet.'}
            </p>
          </div>
        ) : (
          <div className="overflow-hidden rounded-lg border border-border">
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border bg-muted/50">
                    <th className="px-4 py-3 text-left font-medium text-muted-foreground">Key</th>
                    <th className="px-4 py-3 text-left font-medium text-muted-foreground">Title</th>
                    <th className="hidden px-4 py-3 text-left font-medium text-muted-foreground sm:table-cell">Project</th>
                    <th className="px-4 py-3 text-left font-medium text-muted-foreground">Status</th>
                    <th className="hidden px-4 py-3 text-left font-medium text-muted-foreground md:table-cell">Priority</th>
                    <th className="hidden px-4 py-3 text-left font-medium text-muted-foreground lg:table-cell">Due Date</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {tasks.map((task) => {
                    const sc = STATUS_CONFIG[task.status];
                    const overdue = isOverdue(task);
                    return (
                      <tr
                        key={task.id}
                        className="group transition-colors hover:bg-muted/30"
                      >
                        <td className="whitespace-nowrap px-4 py-3">
                          <span className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs text-muted-foreground">
                            {task.task_key}
                          </span>
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center gap-2">
                            <span className="font-medium text-foreground">{task.title}</span>
                            {overdue && (
                              <span className="inline-flex items-center gap-0.5 rounded-full bg-red-500/10 px-1.5 py-0.5 text-[10px] font-semibold text-red-500">
                                OVERDUE
                              </span>
                            )}
                          </div>
                        </td>
                        <td className="hidden whitespace-nowrap px-4 py-3 text-muted-foreground sm:table-cell">
                          {projectMap.get(task.project_id) ?? '—'}
                        </td>
                        <td className="whitespace-nowrap px-4 py-3">
                          <span className="inline-flex items-center gap-1.5">
                            <span className={cn('h-2 w-2 rounded-full', sc.color)} />
                            <span className="text-xs">{sc.label}</span>
                          </span>
                        </td>
                        <td className="hidden whitespace-nowrap px-4 py-3 md:table-cell">
                          <span className={cn('rounded-full px-2 py-0.5 text-xs font-medium', PRIORITY_COLORS[task.priority])}>
                            {task.priority}
                          </span>
                        </td>
                        <td className={cn('hidden whitespace-nowrap px-4 py-3 text-xs lg:table-cell', overdue ? 'text-red-500 font-medium' : 'text-muted-foreground')}>
                          {formatDate(task.due_date)}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>

            {total > pageSize && (
              <div className="border-t border-border px-4 py-3">
                <Pagination
                  page={page}
                  pageSize={pageSize}
                  total={total}
                  onPageChange={setPage}
                  onPageSizeChange={(s) => { setPageSize(s); setPage(1); }}
                  pageSizeOptions={[10, 20, 50]}
                />
              </div>
            )}
          </div>
        )}
      </div>
    </PageShell>
  );
}

function StatCard({ label, value, icon: Icon, color }: { label: string; value: number; icon: typeof Clock; color: string }) {
  return (
    <div className="rounded-lg border border-border bg-card p-4 transition-shadow hover:shadow-sm">
      <div className="flex items-center justify-between">
        <Icon className={cn('h-4 w-4', color)} />
        <span className={cn('text-2xl font-bold tabular-nums', color)}>{value}</span>
      </div>
      <p className="mt-1 text-xs text-muted-foreground">{label}</p>
    </div>
  );
}

function FilterSelect({ label, value, onChange, options }: {
  label: string; value: string; onChange: (v: string) => void;
  options: { value: string; label: string }[];
}) {
  return (
    <div className="space-y-1">
      <label className="text-xs font-medium text-muted-foreground">{label}</label>
      <select
        value={value}
        onChange={e => onChange(e.target.value)}
        className="w-full rounded-md border border-border bg-background px-2 py-1.5 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
      >
        {options.map(o => <option key={o.value} value={o.value}>{o.label}</option>)}
      </select>
    </div>
  );
}
