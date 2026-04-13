import { useState, useMemo } from 'react';
import { UserPlus, Users, Trash2, LogOut } from 'lucide-react';
import type { ProjectMember, WorkspaceMember, WorkspaceRole, Team } from '@/types';
import { roleLabels } from '@utils/format';
import { cn } from '@utils/cn';
import { Badge } from '@components/ui/Badge';

interface Props {
  projectMembers: ProjectMember[];
  workspaceMembers: WorkspaceMember[];
  orgTeams: Team[];
  canManage: boolean;
  currentUserId: string;
  onAdd: (userId: string, role: WorkspaceRole) => void;
  onRemove: (userId: string) => void;
  onChangeRole: (userId: string, role: WorkspaceRole) => void;
  onAddTeam: (teamId: string, defaultRole: string) => void;
  onLeave: () => void;
}

const assignableRolesList: WorkspaceRole[] = ['admin', 'manager', 'lead', 'member', 'guest', 'viewer'];

export function ProjectMemberList({
  projectMembers, workspaceMembers, orgTeams, canManage, currentUserId,
  onAdd, onRemove, onChangeRole, onAddTeam, onLeave,
}: Props) {
  const [showAdd, setShowAdd] = useState(false);
  const [addMode, setAddMode] = useState<'member' | 'team'>('member');
  const [selectedUserId, setSelectedUserId] = useState('');
  const [selectedTeamId, setSelectedTeamId] = useState('');
  const [role, setRole] = useState<WorkspaceRole | ''>('');
  const [editingRole, setEditingRole] = useState<string | null>(null);

  const projectMemberIds = useMemo(() => new Set(projectMembers.map(m => m.user_id)), [projectMembers]);

  const availableWsMembers = useMemo(
    () => workspaceMembers.filter(m => !projectMemberIds.has(m.user_id)),
    [workspaceMembers, projectMemberIds],
  );

  const handleAdd = () => {
    if (!selectedUserId || !role) return;
    onAdd(selectedUserId, role as WorkspaceRole);
    setSelectedUserId('');
    setRole('');
    setShowAdd(false);
  };

  const handleAddTeam = () => {
    if (!selectedTeamId || !role) return;
    onAddTeam(selectedTeamId, role);
    setSelectedTeamId('');
    setRole('');
    setShowAdd(false);
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium text-foreground">Project Members ({projectMembers.length})</h3>
        <div className="flex gap-2">
          <button
            onClick={onLeave}
            className="inline-flex items-center gap-1 rounded-md border border-border px-2 py-1 text-xs text-muted-foreground hover:text-destructive hover:border-destructive/50 transition-colors"
          >
            <LogOut className="h-3 w-3" /> Leave
          </button>
          {canManage && (
            <button
              onClick={() => setShowAdd(!showAdd)}
              className="inline-flex items-center gap-1 rounded-md bg-primary/10 px-2.5 py-1 text-xs font-medium text-primary hover:bg-primary/20"
            >
              <UserPlus className="h-3 w-3" /> Add
            </button>
          )}
        </div>
      </div>

      {showAdd && (
        <div className="rounded-lg border border-border bg-muted/30 p-3 space-y-3">
          <div className="flex gap-1">
            <button
              onClick={() => setAddMode('member')}
              className={cn(
                'px-3 py-1 text-xs font-medium rounded-md transition-colors',
                addMode === 'member' ? 'bg-primary text-primary-foreground' : 'bg-background text-muted-foreground hover:text-foreground',
              )}
            >
              <UserPlus className="inline h-3 w-3 mr-1" />Member
            </button>
            {orgTeams.length > 0 && (
              <button
                onClick={() => setAddMode('team')}
                className={cn(
                  'px-3 py-1 text-xs font-medium rounded-md transition-colors',
                  addMode === 'team' ? 'bg-primary text-primary-foreground' : 'bg-background text-muted-foreground hover:text-foreground',
                )}
              >
                <Users className="inline h-3 w-3 mr-1" />Team
              </button>
            )}
          </div>

          {addMode === 'member' && (
            <>
              <select
                value={selectedUserId}
                onChange={e => setSelectedUserId(e.target.value)}
                className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              >
                <option value="">Select workspace member...</option>
                {availableWsMembers.map(m => (
                  <option key={m.user_id} value={m.user_id}>{m.user_name} ({m.user_email})</option>
                ))}
              </select>
              <div className="flex gap-2">
                <select
                  value={role}
                  onChange={e => setRole(e.target.value as WorkspaceRole)}
                  className="flex-1 rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                >
                  <option value="">Select role...</option>
                  {assignableRolesList.map(r => (
                    <option key={r} value={r}>{roleLabels[r]}</option>
                  ))}
                </select>
                <button
                  onClick={handleAdd}
                  disabled={!selectedUserId || !role}
                  className="rounded-lg bg-primary px-3 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                >
                  Add
                </button>
              </div>
            </>
          )}

          {addMode === 'team' && (
            <>
              <select
                value={selectedTeamId}
                onChange={e => setSelectedTeamId(e.target.value)}
                className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              >
                <option value="">Select team...</option>
                {orgTeams.map(t => (
                  <option key={t.id} value={t.id}>{t.name}</option>
                ))}
              </select>
              <div className="flex gap-2">
                <select
                  value={role}
                  onChange={e => setRole(e.target.value as WorkspaceRole)}
                  className="flex-1 rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                >
                  <option value="">Default role...</option>
                  {assignableRolesList.map(r => (
                    <option key={r} value={r}>{roleLabels[r]}</option>
                  ))}
                </select>
                <button
                  onClick={handleAddTeam}
                  disabled={!selectedTeamId || !role}
                  className="rounded-lg bg-primary px-3 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                >
                  Add Team
                </button>
              </div>
            </>
          )}
        </div>
      )}

      <div className="space-y-2">
        {projectMembers.map(m => (
          <div key={m.id} className="flex items-center justify-between rounded-lg border border-border bg-card p-3">
            <div className="flex items-center gap-2.5">
              <div className="flex h-7 w-7 items-center justify-center rounded-full bg-primary/10 text-xs font-medium text-primary">
                {m.user_name.charAt(0).toUpperCase()}
              </div>
              <div>
                <p className="text-sm font-medium text-foreground">
                  {m.user_name}
                  {m.user_id === currentUserId && <span className="text-muted-foreground"> (you)</span>}
                </p>
                <p className="text-xs text-muted-foreground">{m.user_email}</p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              {canManage && editingRole === m.user_id ? (
                <select
                  value={m.role}
                  onChange={e => {
                    onChangeRole(m.user_id, e.target.value as WorkspaceRole);
                    setEditingRole(null);
                  }}
                  onBlur={() => setEditingRole(null)}
                  autoFocus
                  className="rounded border border-border bg-background px-2 py-0.5 text-xs focus:outline-none focus:ring-1 focus:ring-ring"
                >
                  {assignableRolesList.map(r => (
                    <option key={r} value={r}>{roleLabels[r]}</option>
                  ))}
                </select>
              ) : (
                <button
                  onClick={() => canManage && m.role !== 'owner' && setEditingRole(m.user_id)}
                  className={cn(canManage && m.role !== 'owner' && 'cursor-pointer hover:opacity-70')}
                >
                  <Badge variant={m.role}>{roleLabels[m.role]}</Badge>
                </button>
              )}
              {canManage && m.user_id !== currentUserId && m.role !== 'owner' && (
                <button
                  onClick={() => onRemove(m.user_id)}
                  className="text-destructive hover:text-destructive/80 transition-colors"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              )}
            </div>
          </div>
        ))}
        {projectMembers.length === 0 && (
          <p className="text-center text-sm text-muted-foreground py-6">No project members yet. Add workspace members to this project.</p>
        )}
      </div>
    </div>
  );
}
