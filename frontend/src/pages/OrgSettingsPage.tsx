import { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  Settings, Users, UserPlus, Mail, Shield, Loader2, Trash2, ChevronLeft,
  BarChart3, CheckCircle2, Clock, AlertTriangle, CircleDot, ChevronDown, ChevronUp,
} from 'lucide-react';
import { PageShell } from '@components/layout/PageShell';
import { Pagination } from '@components/ui/Pagination';
import { useAuth } from '@hooks/useAuth';
import { useOrgMembers } from '@hooks/useOrganizations';
import { useOrgInvitations, useSendOrgInvite } from '@hooks/useOrgInvitations';
import { useTeams, useCreateTeam, useTeamDetail, useAddTeamMember, useRemoveTeamMember, useDeleteTeam } from '@hooks/useTeams';
import { useRemoveOrgMember } from '@hooks/useWorkspaces';
import { useOrgMemberStats, useOrgMemberTasks } from '@hooks/useOrgDashboard';
import type { OrgRole, Team, Task } from '@/types';
import { cn } from '@utils/cn';

type Tab = 'members' | 'teams' | 'invitations' | 'stats';

const orgRoleLabels: Record<OrgRole, string> = {
  owner: 'Owner', admin: 'Admin', manager: 'Manager', member: 'Member',
};

const canManageOrg = (role?: OrgRole) => role === 'owner' || role === 'admin';
const canManageTeams = (role?: OrgRole) => role === 'owner' || role === 'admin' || role === 'manager';

export function OrgSettingsPage() {
  const { orgSlug } = useParams();
  const { activeOrg, user } = useAuth();
  const orgId = activeOrg?.id || '';
  const [tab, setTab] = useState<Tab>('members');

  const { data: members } = useOrgMembers(orgId);
  const currentMember = members?.find(m => m.user_id === user?.id);
  const myOrgRole = currentMember?.role;

  return (
    <PageShell>
      <div className="mx-auto max-w-4xl">
        <div className="mb-6">
          <Link to={`/org/${orgSlug}/workspaces`} className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-3">
            <ChevronLeft className="h-4 w-4" /> Back to workspaces
          </Link>
          <div className="flex items-center gap-3">
            <Settings className="h-6 w-6 text-primary" />
            <h1 className="text-2xl font-bold text-foreground">Organization Settings</h1>
          </div>
          <p className="mt-1 text-muted-foreground">{activeOrg?.name}</p>
        </div>

        <div className="flex gap-1 mb-6 border-b border-border overflow-x-auto">
          {(['members', 'stats', 'teams', 'invitations'] as Tab[]).map(t => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={cn(
                'whitespace-nowrap px-4 py-2.5 text-sm font-medium transition-colors border-b-2 -mb-px',
                tab === t
                  ? 'border-primary text-foreground'
                  : 'border-transparent text-muted-foreground hover:text-foreground',
                (t === 'teams' || t === 'stats') && !canManageTeams(myOrgRole) && 'hidden',
              )}
            >
              {t === 'members' && 'Members'}
              {t === 'stats' && 'Member Stats'}
              {t === 'teams' && 'Teams'}
              {t === 'invitations' && 'Invitations'}
            </button>
          ))}
        </div>

        {tab === 'members' && <MembersTab orgId={orgId} myRole={myOrgRole} />}
        {tab === 'stats' && <MemberStatsTab orgId={orgId} />}
        {tab === 'teams' && <TeamsTab orgId={orgId} myRole={myOrgRole} />}
        {tab === 'invitations' && <InvitationsTab orgId={orgId} myRole={myOrgRole} />}
      </div>
    </PageShell>
  );
}

function MembersTab({ orgId, myRole }: { orgId: string; myRole?: OrgRole }) {
  const { data: members, isLoading } = useOrgMembers(orgId);
  const removeMember = useRemoveOrgMember(orgId);

  if (isLoading) return <Loader2 className="mx-auto h-8 w-8 animate-spin text-muted-foreground" />;

  return (
    <div className="space-y-3">
      {members?.map(m => (
        <div key={m.id} className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
          <div>
            <p className="font-medium text-foreground">{m.user_name}</p>
            <p className="text-sm text-muted-foreground">{m.user_email}</p>
          </div>
          <div className="flex items-center gap-3">
            <span className={cn(
              'inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium',
              m.role === 'owner' ? 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400'
                : m.role === 'admin' ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
                : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
            )}>
              <Shield className="h-3 w-3" />
              {orgRoleLabels[m.role]}
            </span>
            {canManageOrg(myRole) && m.role !== 'owner' && (
              <button
                onClick={() => removeMember.mutate(m.user_id)}
                className="text-destructive hover:text-destructive/80 transition-colors"
                title="Remove member"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            )}
          </div>
        </div>
      ))}
      {(!members || members.length === 0) && (
        <p className="text-center text-muted-foreground py-8">No members yet</p>
      )}
    </div>
  );
}

function TeamsTab({ orgId, myRole }: { orgId: string; myRole?: OrgRole }) {
  const { data: teams, isLoading } = useTeams(orgId);
  const createTeam = useCreateTeam(orgId);
  const deleteTeam = useDeleteTeam(orgId);
  const [newTeamName, setNewTeamName] = useState('');
  const [selectedTeam, setSelectedTeam] = useState<Team | null>(null);

  if (isLoading) return <Loader2 className="mx-auto h-8 w-8 animate-spin text-muted-foreground" />;

  const handleCreate = () => {
    if (!newTeamName.trim()) return;
    createTeam.mutate(newTeamName.trim(), {
      onSuccess: () => setNewTeamName(''),
    });
  };

  if (selectedTeam) {
    return <TeamDetail orgId={orgId} team={selectedTeam} myRole={myRole} onBack={() => setSelectedTeam(null)} />;
  }

  return (
    <div className="space-y-4">
      {canManageTeams(myRole) && (
        <div className="flex gap-2">
          <input
            value={newTeamName}
            onChange={e => setNewTeamName(e.target.value)}
            placeholder="New team name"
            className="flex-1 rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            onKeyDown={e => e.key === 'Enter' && handleCreate()}
          />
          <button
            onClick={handleCreate}
            disabled={!newTeamName.trim() || createTeam.isPending}
            className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
          >
            {createTeam.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <UserPlus className="h-4 w-4" />}
            Create
          </button>
        </div>
      )}

      <div className="space-y-3">
        {teams?.map(t => (
          <div key={t.id} className="flex items-center justify-between rounded-lg border border-border bg-card p-4">
            <button onClick={() => setSelectedTeam(t)} className="text-left">
              <p className="font-medium text-foreground hover:text-primary transition-colors">{t.name}</p>
            </button>
            {canManageTeams(myRole) && (
              <button
                onClick={() => deleteTeam.mutate(t.id)}
                className="text-destructive hover:text-destructive/80 transition-colors"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            )}
          </div>
        ))}
        {(!teams || teams.length === 0) && (
          <p className="text-center text-muted-foreground py-8">No teams created yet</p>
        )}
      </div>
    </div>
  );
}

function TeamDetail({ orgId, team, myRole, onBack }: {
  orgId: string; team: Team; myRole?: OrgRole; onBack: () => void;
}) {
  const { data, isLoading } = useTeamDetail(orgId, team.id);
  const { data: orgMembers } = useOrgMembers(orgId);
  const addMember = useAddTeamMember(orgId, team.id);
  const removeMember = useRemoveTeamMember(orgId, team.id);
  const [addUserId, setAddUserId] = useState('');

  const teamMembers = data?.members || [];
  const teamMemberIds = new Set(teamMembers.map(m => m.user_id));
  const availableMembers = orgMembers?.filter(m => !teamMemberIds.has(m.user_id)) || [];

  const handleAdd = () => {
    if (!addUserId) return;
    addMember.mutate(addUserId, { onSuccess: () => setAddUserId('') });
  };

  return (
    <div>
      <button onClick={onBack} className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-4">
        <ChevronLeft className="h-4 w-4" /> Back to teams
      </button>
      <h2 className="text-lg font-semibold text-foreground mb-4">
        <Users className="inline h-5 w-5 mr-2 text-primary" />
        {team.name}
      </h2>

      {canManageTeams(myRole) && availableMembers.length > 0 && (
        <div className="flex gap-2 mb-4">
          <select
            value={addUserId}
            onChange={e => setAddUserId(e.target.value)}
            className="flex-1 rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
          >
            <option value="">Select org member to add...</option>
            {availableMembers.map(m => (
              <option key={m.user_id} value={m.user_id}>{m.user_name} ({m.user_email})</option>
            ))}
          </select>
          <button
            onClick={handleAdd}
            disabled={!addUserId || addMember.isPending}
            className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
          >
            Add
          </button>
        </div>
      )}

      {isLoading ? (
        <Loader2 className="mx-auto h-8 w-8 animate-spin text-muted-foreground" />
      ) : (
        <div className="space-y-2">
          {teamMembers.map(m => (
            <div key={m.id} className="flex items-center justify-between rounded-lg border border-border bg-card p-3">
              <div>
                <p className="text-sm font-medium text-foreground">{m.user_name}</p>
                <p className="text-xs text-muted-foreground">{m.user_email}</p>
              </div>
              {canManageTeams(myRole) && (
                <button
                  onClick={() => removeMember.mutate(m.user_id)}
                  className="text-destructive hover:text-destructive/80 transition-colors"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              )}
            </div>
          ))}
          {teamMembers.length === 0 && (
            <p className="text-center text-sm text-muted-foreground py-4">No team members yet</p>
          )}
        </div>
      )}
    </div>
  );
}

function InvitationsTab({ orgId, myRole }: { orgId: string; myRole?: OrgRole }) {
  const { data: invitations, isLoading } = useOrgInvitations(orgId);
  const sendInvite = useSendOrgInvite(orgId);
  const [email, setEmail] = useState('');
  const [role, setRole] = useState<OrgRole>('member');

  const handleSend = () => {
    if (!email.trim()) return;
    sendInvite.mutate({ email: email.trim(), role }, {
      onSuccess: () => { setEmail(''); setRole('member'); },
    });
  };

  if (isLoading) return <Loader2 className="mx-auto h-8 w-8 animate-spin text-muted-foreground" />;

  return (
    <div className="space-y-4">
      {canManageOrg(myRole) && (
        <div className="rounded-xl border border-border bg-card p-4 space-y-3">
          <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
            <Mail className="h-4 w-4 text-primary" /> Invite to Organization
          </h3>
          <div className="flex gap-2">
            <input
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              placeholder="user@example.com"
              className="flex-1 rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            />
            <select
              value={role}
              onChange={e => setRole(e.target.value as OrgRole)}
              className="rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            >
              <option value="member">Member</option>
              <option value="manager">Manager</option>
              <option value="admin">Admin</option>
            </select>
            <button
              onClick={handleSend}
              disabled={!email.trim() || sendInvite.isPending}
              className="inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50"
            >
              {sendInvite.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <UserPlus className="h-4 w-4" />}
              Send
            </button>
          </div>
        </div>
      )}

      <div className="space-y-2">
        {invitations?.map(inv => (
          <div key={inv.id} className="flex items-center justify-between rounded-lg border border-border bg-card p-3">
            <div>
              <p className="text-sm font-medium text-foreground">{inv.invitee_email}</p>
              <p className="text-xs text-muted-foreground">
                Role: {orgRoleLabels[inv.role]} &middot; Invited by {inv.inviter_name}
              </p>
            </div>
            <span className={cn(
              'rounded-full px-2.5 py-0.5 text-xs font-medium',
              inv.status === 'pending' ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
                : inv.status === 'accepted' ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
            )}>
              {inv.status}
            </span>
          </div>
        ))}
        {(!invitations || invitations.length === 0) && (
          <p className="text-center text-muted-foreground py-8">No invitations sent yet</p>
        )}
      </div>
    </div>
  );
}

const STATUS_COLORS: Record<string, string> = {
  todo: 'bg-slate-500', in_progress: 'bg-blue-500', blocked: 'bg-red-500', done: 'bg-emerald-500',
};
const STATUS_LABELS: Record<string, string> = {
  todo: 'To Do', in_progress: 'In Progress', blocked: 'Blocked', done: 'Done',
};
const PRIORITY_BADGE: Record<string, string> = {
  high: 'text-red-500 bg-red-500/10', medium: 'text-amber-500 bg-amber-500/10', low: 'text-emerald-500 bg-emerald-500/10',
};

function MemberStatsTab({ orgId }: { orgId: string }) {
  const { data: memberStats, isLoading, isError } = useOrgMemberStats(orgId);
  const [expandedUser, setExpandedUser] = useState<string | null>(null);

  if (isLoading) return <Loader2 className="mx-auto h-8 w-8 animate-spin text-muted-foreground" />;
  if (isError) return <p className="text-center text-destructive py-8">Failed to load member stats.</p>;

  const stats = memberStats ?? [];
  if (stats.length === 0) return <p className="text-center text-muted-foreground py-8">No members below your role level in this organization.</p>;

  const totalTasks = stats.reduce((s, m) => s + m.total, 0);
  const totalOverdue = stats.reduce((s, m) => s + m.overdue, 0);
  const totalCompleted = stats.reduce((s, m) => s + m.completed, 0);

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <MiniStat label="Members" value={stats.length} icon={Users} color="text-foreground" />
        <MiniStat label="Total Tasks" value={totalTasks} icon={BarChart3} color="text-blue-500" />
        <MiniStat label="Completed" value={totalCompleted} icon={CheckCircle2} color="text-emerald-500" />
        <MiniStat label="Overdue" value={totalOverdue} icon={AlertTriangle} color="text-red-500" />
      </div>

      <div className="overflow-hidden rounded-lg border border-border">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border bg-muted/50">
              <th className="px-4 py-3 text-left font-medium text-muted-foreground">Member</th>
              <th className="hidden px-4 py-3 text-left font-medium text-muted-foreground sm:table-cell">Role</th>
              <th className="px-4 py-3 text-center font-medium text-muted-foreground">Total</th>
              <th className="hidden px-4 py-3 text-center font-medium text-muted-foreground md:table-cell">
                <CircleDot className="mx-auto h-4 w-4 text-slate-400" title="To Do" />
              </th>
              <th className="hidden px-4 py-3 text-center font-medium text-muted-foreground md:table-cell">
                <Clock className="mx-auto h-4 w-4 text-blue-500" title="In Progress" />
              </th>
              <th className="hidden px-4 py-3 text-center font-medium text-muted-foreground md:table-cell">
                <AlertTriangle className="mx-auto h-4 w-4 text-red-500" title="Blocked" />
              </th>
              <th className="hidden px-4 py-3 text-center font-medium text-muted-foreground md:table-cell">
                <CheckCircle2 className="mx-auto h-4 w-4 text-emerald-500" title="Done" />
              </th>
              <th className="px-4 py-3 text-center font-medium text-muted-foreground">
                <span className="text-xs text-red-500">Overdue</span>
              </th>
              <th className="px-2 py-3" />
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {stats.map(m => (
              <MemberStatsRow
                key={m.user_id}
                member={m}
                orgId={orgId}
                isExpanded={expandedUser === m.user_id}
                onToggle={() => setExpandedUser(expandedUser === m.user_id ? null : m.user_id)}
              />
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function MemberStatsRow({ member: m, orgId, isExpanded, onToggle }: {
  member: { user_id: string; user_name: string; user_email: string; role: string; total: number; by_status: Record<string, number>; overdue: number; completed: number };
  orgId: string; isExpanded: boolean; onToggle: () => void;
}) {
  return (
    <>
      <tr className="group hover:bg-muted/30 transition-colors">
        <td className="px-4 py-3">
          <p className="font-medium text-foreground">{m.user_name}</p>
          <p className="text-xs text-muted-foreground">{m.user_email}</p>
        </td>
        <td className="hidden px-4 py-3 sm:table-cell">
          <span className={cn(
            'inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium',
            m.role === 'admin' ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
              : m.role === 'manager' ? 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-400'
              : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
          )}>
            {m.role}
          </span>
        </td>
        <td className="px-4 py-3 text-center font-semibold tabular-nums">{m.total}</td>
        <td className="hidden px-4 py-3 text-center tabular-nums md:table-cell">{m.by_status.todo ?? 0}</td>
        <td className="hidden px-4 py-3 text-center tabular-nums md:table-cell">{m.by_status.in_progress ?? 0}</td>
        <td className="hidden px-4 py-3 text-center tabular-nums md:table-cell">{m.by_status.blocked ?? 0}</td>
        <td className="hidden px-4 py-3 text-center tabular-nums md:table-cell">{m.by_status.done ?? 0}</td>
        <td className={cn('px-4 py-3 text-center tabular-nums font-medium', m.overdue > 0 ? 'text-red-500' : 'text-muted-foreground')}>
          {m.overdue}
        </td>
        <td className="px-2 py-3">
          <button onClick={onToggle} className="p-1 text-muted-foreground hover:text-foreground transition-colors" title="View tasks">
            {isExpanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
          </button>
        </td>
      </tr>
      {isExpanded && (
        <tr>
          <td colSpan={9} className="bg-muted/20 p-0">
            <MemberTasksExpanded orgId={orgId} userId={m.user_id} userName={m.user_name} />
          </td>
        </tr>
      )}
    </>
  );
}

function MemberTasksExpanded({ orgId, userId, userName }: { orgId: string; userId: string; userName: string }) {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);
  const [statusFilter, setStatusFilter] = useState('');
  const { data, isLoading } = useOrgMemberTasks(orgId, userId, { status: statusFilter || undefined, page, limit: pageSize });

  const tasks = data?.data ?? [];
  const total = data?.meta?.total ?? 0;

  const isOverdue = (t: Task) => t.due_date && t.status !== 'done' && new Date(t.due_date) < new Date(new Date().toDateString());
  const formatDate = (d: string | null) => {
    if (!d) return '—';
    return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  return (
    <div className="px-6 py-4 space-y-3">
      <div className="flex items-center justify-between">
        <p className="text-sm font-medium text-foreground">{userName}'s Tasks</p>
        <select
          value={statusFilter}
          onChange={e => { setStatusFilter(e.target.value); setPage(1); }}
          className="rounded-md border border-border bg-background px-2 py-1 text-xs text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        >
          <option value="">All statuses</option>
          <option value="todo">To Do</option>
          <option value="in_progress">In Progress</option>
          <option value="blocked">Blocked</option>
          <option value="done">Done</option>
        </select>
      </div>

      {isLoading ? (
        <Loader2 className="mx-auto h-5 w-5 animate-spin text-muted-foreground" />
      ) : tasks.length === 0 ? (
        <p className="text-center text-xs text-muted-foreground py-3">No tasks found.</p>
      ) : (
        <div className="space-y-1.5">
          {tasks.map(t => {
            const overdue = isOverdue(t);
            return (
              <div key={t.id} className="flex items-center gap-3 rounded-md border border-border bg-card px-3 py-2">
                <span className={cn('h-2 w-2 shrink-0 rounded-full', STATUS_COLORS[t.status])} title={STATUS_LABELS[t.status]} />
                <span className="rounded bg-muted px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">{t.task_key}</span>
                <span className="flex-1 truncate text-sm text-foreground">{t.title}</span>
                <span className={cn('rounded-full px-1.5 py-0.5 text-[10px] font-medium', PRIORITY_BADGE[t.priority])}>
                  {t.priority}
                </span>
                <span className={cn('text-xs tabular-nums', overdue ? 'text-red-500 font-medium' : 'text-muted-foreground')}>
                  {formatDate(t.due_date)}
                </span>
              </div>
            );
          })}
        </div>
      )}

      {total > pageSize && (
        <Pagination page={page} pageSize={pageSize} total={total} onPageChange={setPage}
          onPageSizeChange={s => { setPageSize(s); setPage(1); }} pageSizeOptions={[5, 10, 20]} />
      )}
    </div>
  );
}

function MiniStat({ label, value, icon: Icon, color }: { label: string; value: number; icon: typeof Users; color: string }) {
  return (
    <div className="rounded-lg border border-border bg-card p-3">
      <div className="flex items-center justify-between">
        <Icon className={cn('h-4 w-4', color)} />
        <span className={cn('text-xl font-bold tabular-nums', color)}>{value}</span>
      </div>
      <p className="mt-0.5 text-xs text-muted-foreground">{label}</p>
    </div>
  );
}
