import { useState, useMemo, useCallback, useRef } from 'react';
import { useParams } from 'react-router-dom';
import { Plus, Loader2, Filter, Search, MessageSquare, Settings, X, ChevronDown, ChevronUp } from 'lucide-react';
import { PageShell } from '@components/layout/PageShell';
import { Breadcrumbs } from '@components/layout/Breadcrumbs';
import { Select } from '@components/ui/Select';
import { TaskBoard } from '@features/tasks/TaskBoard';
import { TaskModal } from '@features/tasks/TaskModal';
import { BlockedReasonDialog } from '@features/tasks/BlockedReasonDialog';
import { CustomFieldsManager } from '@features/projects/CustomFieldsManager';
import { WIPLimitsManager } from '@features/projects/WIPLimitsManager';
import { CommentThread } from '@features/comments/CommentThread';
import { ProjectMemberList } from '@features/projects/ProjectMemberList';
import { useProject } from '@hooks/useProjects';
import { useWorkspace, useWorkspaceMembers } from '@hooks/useWorkspaces';
import { useTasks, useCreateTask, useUpdateTask, useDeleteTask } from '@hooks/useTasks';
import { useCustomFieldDefinitions, useFieldValues, useSetFieldValue, useWIPLimits } from '@hooks/useCustomFields';
import { useProjectMembers, useAddProjectMember, useRemoveProjectMember, useUpdateProjectMemberRole, useLeaveProject } from '@hooks/useProjectMembers';
import { useTeams, useAddTeamToProject } from '@hooks/useTeams';
import { useAuth } from '@hooks/useAuth';
import type { Task, TaskStatus, WIPLimit, WorkspaceRole } from '@/types';
import { cn } from '@utils/cn';
import { statusLabels, rolePower } from '@utils/format';
import { toast } from 'sonner';

type SettingsTab = 'members' | 'custom-fields' | 'wip-limits';

const priorityFilterOptions = [
  { value: '', label: 'All priorities' },
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
];

export function ProjectDetailPage() {
  const { orgSlug, wid, pid } = useParams<{ orgSlug: string; wid: string; pid: string }>();
  const workspaceId = wid!;
  const projectId = pid!;
  const { user, activeOrg } = useAuth();
  const orgId = activeOrg?.id ?? '';

  const [modalOpen, setModalOpen] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const [modalKey, setModalKey] = useState(0);
  const [filterAssignee, setFilterAssignee] = useState('');
  const [filterPriority, setFilterPriority] = useState('');
  const [filterSearch, setFilterSearch] = useState('');
  const [discussionOpen, setDiscussionOpen] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [activeSettingsTab, setActiveSettingsTab] = useState<SettingsTab>('members');

  const [blockedDialogOpen, setBlockedDialogOpen] = useState(false);
  const pendingBlockedTaskId = useRef<string | null>(null);
  const pendingBlockedSource = useRef<'board' | 'modal'>('board');
  const { data: workspace } = useWorkspace(orgId, workspaceId);
  const { data: project, isLoading: projLoading, isError: projError } = useProject(projectId);
  const { data: members = [] } = useWorkspaceMembers(orgId, workspaceId);
  const { data: tasksData, isLoading: tasksLoading, isError: tasksError } = useTasks(projectId, {
    assignee: filterAssignee || undefined,
    priority: filterPriority || undefined,
    search: filterSearch || undefined,
  });

  const { data: customFields = [] } = useCustomFieldDefinitions(projectId);
  const { data: wipLimits = [] } = useWIPLimits(projectId);
  const { data: editingFieldValues = [] } = useFieldValues(editingTask?.id ?? '');
  const setFieldValue = useSetFieldValue(editingTask?.id ?? '');

  const { data: projectMembers = [] } = useProjectMembers(orgId, workspaceId, projectId);
  const addProjectMember = useAddProjectMember(orgId, workspaceId, projectId);
  const removeProjectMember = useRemoveProjectMember(orgId, workspaceId, projectId);
  const updateProjectMemberRole = useUpdateProjectMemberRole(orgId, workspaceId, projectId);
  const leaveProject = useLeaveProject(orgId, workspaceId, projectId);
  const { data: orgTeams } = useTeams(orgId);
  const addTeamToProject = useAddTeamToProject(orgId, workspaceId, projectId);

  const createTask = useCreateTask(projectId);
  const updateTask = useUpdateTask(projectId);
  const deleteTask = useDeleteTask(projectId);

  const currentMember = members.find((m) => m.user_id === user?.id);
  const currentRole = (currentMember?.role ?? 'member') as WorkspaceRole;
  const canManageSettings = rolePower[currentRole] >= rolePower.admin;

  const tasks = useMemo(() => tasksData?.data ?? [], [tasksData]);

  const wipLimitMap = useMemo(
    () => new Map<TaskStatus, WIPLimit>(wipLimits.map((l) => [l.status, l])),
    [wipLimits],
  );

  const taskCountByStatus = useMemo(() => {
    const counts: Record<TaskStatus, number> = { todo: 0, in_progress: 0, blocked: 0, done: 0 };
    for (const t of tasks) counts[t.status]++;
    return counts;
  }, [tasks]);

  const customFieldValuesMap = useMemo(() => {
    const map: Record<string, string> = {};
    for (const v of editingFieldValues) map[v.field_id] = v.value;
    return map;
  }, [editingFieldValues]);

  const handleStatusChange = useCallback(
    (taskId: string, newStatus: TaskStatus) => {
      if (newStatus === 'blocked') {
        pendingBlockedTaskId.current = taskId;
        pendingBlockedSource.current = 'board';
        setBlockedDialogOpen(true);
        return;
      }

      const limit = wipLimitMap.get(newStatus);
      if (limit && taskCountByStatus[newStatus] >= limit.max_tasks) {
        toast.warning('WIP limit reached', {
          description: `"${statusLabels[newStatus]}" column is at capacity (${limit.max_tasks} max). Move or complete existing tasks first.`,
        });
        return;
      }
      updateTask.mutate({ taskId, body: { status: newStatus } });
    },
    [wipLimitMap, taskCountByStatus, updateTask],
  );

  const handleBlockedConfirm = useCallback(
    (blockedByTask: string, blockedReason: string) => {
      const taskId = pendingBlockedTaskId.current;
      if (!taskId) return;

      updateTask.mutate(
        { taskId, body: { status: 'blocked', blocked_by_task: blockedByTask, blocked_reason: blockedReason } },
        {
          onSuccess: () => {
            if (pendingBlockedSource.current === 'modal') {
              setModalOpen(false);
              setEditingTask(null);
            }
          },
        },
      );
      setBlockedDialogOpen(false);
      pendingBlockedTaskId.current = null;
    },
    [updateTask],
  );

  const handleTaskClick = useCallback((task: Task) => {
    setEditingTask(task);
    setModalOpen(true);
  }, []);

  const handleSave = (data: {
    title: string;
    description: string;
    status?: string;
    priority: string;
    assignee_id: string;
    start_date: string;
    due_date: string;
    custom_fields?: Record<string, string>;
  }) => {
    const { custom_fields, start_date, due_date, assignee_id, ...rest } = data;
    const taskData = {
      ...rest,
      assignee_id: assignee_id || undefined,
      start_date: start_date || undefined,
      due_date: due_date || undefined,
    };

    if (editingTask && taskData.status === 'blocked' && editingTask.status !== 'blocked') {
      pendingBlockedTaskId.current = editingTask.id;
      pendingBlockedSource.current = 'modal';
      setBlockedDialogOpen(true);
      return;
    }

    const onFieldsSave = () => {
      if (custom_fields && editingTask) {
        for (const [fieldId, value] of Object.entries(custom_fields)) {
          if (value !== customFieldValuesMap[fieldId]) {
            setFieldValue.mutate({ fieldId, value });
          }
        }
      }
    };

    if (editingTask) {
      updateTask.mutate(
        { taskId: editingTask.id, body: taskData },
        {
          onSuccess: () => {
            onFieldsSave();
            setModalOpen(false);
            setEditingTask(null);
          },
        },
      );
    } else {
      createTask.mutate(taskData, { onSuccess: () => setModalOpen(false) });
    }
  };

  const handleDelete = () => {
    if (editingTask) {
      deleteTask.mutate(editingTask.id, {
        onSuccess: () => { setModalOpen(false); setEditingTask(null); },
      });
    }
  };

  const assigneeOptions = useMemo(
    () => [
      { value: '', label: 'All assignees' },
      ...members.map((m) => ({ value: m.user_id, label: m.user_name })),
    ],
    [members],
  );

  if (projLoading) {
    return (
      <PageShell>
        <div className="flex justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      </PageShell>
    );
  }

  if (projError) {
    return (
      <PageShell>
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-6 text-center">
          <p className="text-destructive font-medium">Failed to load project</p>
          <p className="mt-1 text-sm text-muted-foreground">The project may have been deleted or you don't have access.</p>
        </div>
      </PageShell>
    );
  }

  return (
    <PageShell>
      <Breadcrumbs items={[
        { label: 'Workspaces', href: `/org/${orgSlug}/workspaces` },
        { label: workspace?.name ?? '...', href: `/org/${orgSlug}/workspaces/${workspaceId}` },
        { label: project?.name ?? '' },
      ]} />

      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-foreground">{project?.name}</h1>
          {project?.description && <p className="mt-1 text-muted-foreground text-sm">{project.description}</p>}
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setSettingsOpen(true)}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <Settings className="h-4 w-4" /> Settings
          </button>
          <button
            onClick={() => { setEditingTask(null); setModalKey((k) => k + 1); setModalOpen(true); }}
            className={cn('inline-flex items-center gap-1 rounded-lg bg-primary px-3 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90')}
          >
            <Plus className="h-4 w-4" /> Add Task
          </button>
        </div>
      </div>

      {/* Filter Bar */}
      <div className="mb-6 flex flex-wrap items-center gap-3">
        <div className="flex items-center gap-2 rounded-lg border border-border bg-card px-3 py-1.5">
          <Search className="h-3.5 w-3.5 text-muted-foreground" />
          <input
            type="text"
            value={filterSearch}
            onChange={(e) => setFilterSearch(e.target.value)}
            placeholder="Search tasks…"
            className="w-40 bg-transparent text-sm text-foreground placeholder:text-muted-foreground focus:outline-none"
          />
        </div>

        <div className="flex items-center gap-2">
          <Filter className="h-3.5 w-3.5 text-muted-foreground" />
          <Select
            value={filterPriority}
            onChange={setFilterPriority}
            options={priorityFilterOptions}
            className="w-40"
          />
        </div>

        <Select
          value={filterAssignee}
          onChange={setFilterAssignee}
          options={assigneeOptions}
          className="w-44"
        />
      </div>

      {tasksLoading ? (
        <div className="flex justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : tasksError ? (
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-6 text-center">
          <p className="text-destructive text-sm">Failed to load tasks. Please try again.</p>
        </div>
      ) : (
        <TaskBoard
          tasks={tasks}
          members={members}
          onStatusChange={handleStatusChange}
          onTaskClick={handleTaskClick}
          wipLimits={wipLimitMap}
        />
      )}

      {/* Project Discussion */}
      <div className="mt-10">
        <button
          type="button"
          onClick={() => setDiscussionOpen((o) => !o)}
          className="flex w-full items-center justify-between rounded-lg border border-border bg-card px-4 py-3 text-sm font-semibold text-foreground transition-colors hover:bg-muted"
        >
          <span className="flex items-center gap-2">
            <MessageSquare className="h-4 w-4" />
            Discussion
          </span>
          {discussionOpen ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
        </button>

        {discussionOpen && (
          <div className="mt-3 rounded-xl border border-border bg-card p-5">
            <CommentThread entityType="project" entityId={projectId} currentUserId={user?.id ?? ''} />
          </div>
        )}
      </div>

      <TaskModal
        key={editingTask?.id ?? `new-${modalKey}`}
        open={modalOpen}
        task={editingTask}
        members={members}
        customFields={customFields}
        customFieldValues={customFieldValuesMap}
        onClose={() => { setModalOpen(false); setEditingTask(null); }}
        onSave={handleSave}
        onStatusChange={handleStatusChange}
        onDelete={editingTask ? handleDelete : undefined}
        isSaving={createTask.isPending || updateTask.isPending}
      >
        {editingTask && (
          <div className="mt-6 border-t border-border pt-6">
            <CommentThread entityType="task" entityId={editingTask.id} currentUserId={user?.id ?? ''} />
          </div>
        )}
      </TaskModal>

      <BlockedReasonDialog
        open={blockedDialogOpen}
        onClose={() => { setBlockedDialogOpen(false); pendingBlockedTaskId.current = null; }}
        onConfirm={handleBlockedConfirm}
      />

      {/* Settings Sidebar */}
      {settingsOpen && (
        <div className="fixed inset-0 z-40" onClick={() => setSettingsOpen(false)}>
          <div className="absolute inset-0 bg-black/30" />
        </div>
      )}
      <div
        className={cn(
          'fixed top-0 right-0 z-50 flex h-full w-full max-w-md flex-col border-l border-border bg-card shadow-xl transition-transform duration-300 ease-in-out',
          settingsOpen ? 'translate-x-0' : 'translate-x-full',
        )}
      >
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 className="flex items-center gap-2 text-lg font-semibold text-foreground">
            <Settings className="h-5 w-5" />
            Project Settings
          </h2>
          <button
            onClick={() => setSettingsOpen(false)}
            className="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="flex border-b border-border">
          {([
            { key: 'members' as const, label: 'Members' },
            { key: 'custom-fields' as const, label: 'Custom Fields' },
            { key: 'wip-limits' as const, label: 'WIP Limits' },
          ]).map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveSettingsTab(tab.key)}
              className={cn(
                'flex-1 px-4 py-3 text-sm font-medium transition-colors',
                activeSettingsTab === tab.key
                  ? 'border-b-2 border-primary text-primary'
                  : 'text-muted-foreground hover:text-foreground',
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>

        <div className="flex-1 overflow-y-auto p-5">
          {activeSettingsTab === 'members' && (
            <ProjectMemberList
              projectMembers={projectMembers}
              workspaceMembers={members}
              orgTeams={orgTeams || []}
              canManage={canManageSettings}
              currentUserId={user?.id ?? ''}
              onAdd={(userId: string, role: WorkspaceRole) => addProjectMember.mutate({ userId, role })}
              onRemove={(userId: string) => removeProjectMember.mutate(userId)}
              onChangeRole={(userId: string, role: WorkspaceRole) => updateProjectMemberRole.mutate({ userId, role })}
              onAddTeam={(teamId: string, defaultRole: string) => addTeamToProject.mutate({ teamId, defaultRole })}
              onLeave={() => leaveProject.mutate()}
            />
          )}
          {activeSettingsTab === 'custom-fields' && (
            <CustomFieldsManager projectId={projectId} canManage={canManageSettings} />
          )}
          {activeSettingsTab === 'wip-limits' && (
            <WIPLimitsManager projectId={projectId} canManage={canManageSettings} />
          )}
        </div>
      </div>
    </PageShell>
  );
}
