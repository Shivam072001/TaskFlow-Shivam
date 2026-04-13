import { BarChart3, CheckCircle2, Clock, AlertTriangle } from 'lucide-react';
import type { WorkspaceStats } from '@/types';

export function StatsBar({ stats }: { stats: WorkspaceStats }) {
  const items = [
    { label: 'Projects', value: stats.project_count, icon: BarChart3, color: 'text-blue-500' },
    { label: 'To Do', value: stats.tasks_by_status?.todo ?? 0, icon: Clock, color: 'text-yellow-500' },
    { label: 'In Progress', value: stats.tasks_by_status?.in_progress ?? 0, icon: BarChart3, color: 'text-primary' },
    { label: 'Done', value: stats.tasks_by_status?.done ?? 0, icon: CheckCircle2, color: 'text-green-500' },
    { label: 'Overdue', value: stats.overdue_count, icon: AlertTriangle, color: 'text-destructive' },
  ];

  return (
    <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-5">
      {items.map((item) => (
        <div key={item.label} className="rounded-xl border border-border bg-card p-4">
          <div className="flex items-center gap-2">
            <item.icon className={`h-4 w-4 ${item.color}`} />
            <span className="text-xs text-muted-foreground">{item.label}</span>
          </div>
          <p className="mt-1 text-2xl font-bold text-foreground">{item.value}</p>
        </div>
      ))}
    </div>
  );
}
