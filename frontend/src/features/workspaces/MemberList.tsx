import { useState, useMemo } from 'react';
import { UserPlus, MoreVertical, Shield, UserMinus, Users, LogOut } from 'lucide-react';
import type { WorkspaceMember, WorkspaceRole, OrgMember, Team } from '@/types';
import { roleLabels, rolePower } from '@utils/format';
import { cn } from '@utils/cn';
import { useAuth } from '@hooks/useAuth';
import type { SelectOption } from '@components/ui/Select';
import { Badge } from '@components/ui/Badge';

interface Props {
  members: WorkspaceMember[];
  orgMembers?: OrgMember[];
  orgTeams?: Team[];
  onDirectAdd?: (userId: string, role: string) => void;
  onAddTeam?: (teamId: string, defaultRole: string) => void;
  onInvite: (email: string, role: string) => void;
  onChangeRole: (userId: string, role: string) => void;
  onRemove: (userId: string) => void;
  onLeave?: () => void;
  currentUserRole: WorkspaceRole;
  isInviting: boolean;
}

export function MemberList(props: Props) {
  const {
    members, orgMembers, orgTeams, onDirectAdd, onAddTeam,
    onChangeRole, onRemove, onLeave, currentUserRole,
  } = props;
  const [showAdd, setShowAdd] = useState(false);
  const [addMode, setAddMode] = useState<'member' | 'team'>('member');
  const [selectedUserId, setSelectedUserId] = useState('');
  const [selectedTeamId, setSelectedTeamId] = useState('');
  const [role, setRole] = useState<WorkspaceRole | ''>('');
  const [openMenu, setOpenMenu] = useState<string | null>(null);
  const { user } = useAuth();

  const myPower = rolePower[currentUserRole];

  const assignableRoles = useMemo(
    () =>
      (Object.keys(rolePower) as WorkspaceRole[])
        .filter((r) => r !== 'owner' && rolePower[r] < myPower)
        .sort((a, b) => rolePower[b] - rolePower[a]),
    [myPower],
  );

  const canManage = myPower >= rolePower['admin'];

  const roleOptions: SelectOption[] = useMemo(
    () => assignableRoles.map((r) => ({ value: r, label: roleLabels[r] })),
    [assignableRoles],
  );

  const canManageMember = (m: WorkspaceMember) =>
    rolePower[m.role] < myPower && m.user_id !== user?.id;

  const memberIds = useMemo(() => new Set(members.map(m => m.user_id)), [members]);
  const availableOrgMembers = useMemo(
    () => (orgMembers || []).filter(m => !memberIds.has(m.user_id)),
    [orgMembers, memberIds],
  );

  const handleDirectAdd = () => {
    if (!selectedUserId || !role) return;
    onDirectAdd?.(selectedUserId, role);
    setSelectedUserId('');
    setRole('');
    setShowAdd(false);
  };

  const handleAddTeam = () => {
    if (!selectedTeamId || !role) return;
    onAddTeam?.(selectedTeamId, role);
    setSelectedTeamId('');
    setRole('');
    setShowAdd(false);
  };

  return (
    <div className="rounded-xl border border-border bg-card">
      <div className="flex items-center justify-between border-b border-border p-4">
        <h3 className="font-semibold text-foreground">Members ({members.length})</h3>
        <div className="flex items-center gap-2">
          {onLeave && (
            <button
              onClick={onLeave}
              className="inline-flex items-center gap-1 rounded-lg border border-border px-2.5 py-1.5 text-xs font-medium text-muted-foreground hover:text-destructive hover:border-destructive/50 transition-colors"
              title="Leave workspace"
            >
              <LogOut className="h-3 w-3" /> Leave
            </button>
          )}
          {canManage && (
            <button
              onClick={() => setShowAdd(!showAdd)}
              className="inline-flex items-center gap-1 rounded-lg bg-primary/10 px-3 py-1.5 text-xs font-medium text-primary hover:bg-primary/20"
            >
              <UserPlus className="h-3 w-3" /> Add
            </button>
          )}
        </div>
      </div>

      {showAdd && (
        <div className="border-b border-border p-4 space-y-3">
          <div className="flex gap-1">
            <button
              onClick={() => setAddMode('member')}
              className={cn(
                'px-3 py-1 text-xs font-medium rounded-md transition-colors',
                addMode === 'member' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground hover:text-foreground',
              )}
            >
              <UserPlus className="inline h-3 w-3 mr-1" />Member
            </button>
            {orgTeams && orgTeams.length > 0 && (
              <button
                onClick={() => setAddMode('team')}
                className={cn(
                  'px-3 py-1 text-xs font-medium rounded-md transition-colors',
                  addMode === 'team' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground hover:text-foreground',
                )}
              >
                <Users className="inline h-3 w-3 mr-1" />Team
              </button>
            )}
          </div>

          {addMode === 'member' && (
            <div className="space-y-2">
              <select
                value={selectedUserId}
                onChange={e => setSelectedUserId(e.target.value)}
                className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
              >
                <option value="">Select org member...</option>
                {availableOrgMembers.map(m => (
                  <option key={m.user_id} value={m.user_id}>{m.user_name} ({m.user_email})</option>
                ))}
              </select>
              <div className="flex gap-2">
                <select
                  value={role}
                  onChange={e => setRole(e.target.value as WorkspaceRole)}
                  className="w-40 rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                >
                  <option value="">Select role...</option>
                  {roleOptions.map(o => (
                    <option key={o.value} value={o.value}>{o.label}</option>
                  ))}
                </select>
                <button
                  onClick={handleDirectAdd}
                  disabled={!selectedUserId || !role}
                  className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                >
                  Add
                </button>
              </div>
            </div>
          )}

          {addMode === 'team' && orgTeams && (
            <div className="space-y-2">
              <select
                value={selectedTeamId}
                onChange={e => setSelectedTeamId(e.target.value)}
                className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
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
                  className="w-40 rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                >
                  <option value="">Default role...</option>
                  {roleOptions.map(o => (
                    <option key={o.value} value={o.value}>{o.label}</option>
                  ))}
                </select>
                <button
                  onClick={handleAddTeam}
                  disabled={!selectedTeamId || !role}
                  className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
                >
                  Add Team
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      <div className="divide-y divide-border">
        {members.map((m) => (
          <div key={m.id} className="flex items-center justify-between p-4">
            <div className="flex items-center gap-3">
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-xs font-medium text-primary">
                {m.user_name.charAt(0).toUpperCase()}
              </div>
              <div>
                <p className="text-sm font-medium text-foreground">
                  {m.user_name}{' '}
                  {m.user_id === user?.id && <span className="text-muted-foreground">(you)</span>}
                </p>
                <p className="text-xs text-muted-foreground">{m.user_email}</p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant={m.role}>{roleLabels[m.role]}</Badge>
              {canManageMember(m) && (
                <div className="relative">
                  <button
                    onClick={() => setOpenMenu(openMenu === m.id ? null : m.id)}
                    className="rounded-md p-1 text-muted-foreground hover:bg-muted hover:text-foreground"
                  >
                    <MoreVertical className="h-4 w-4" />
                  </button>
                  {openMenu === m.id && (
                    <div className="absolute right-0 mt-1 w-48 rounded-lg border border-border bg-card py-1 shadow-lg z-10">
                      <p className="px-3 py-1.5 text-xs font-medium text-muted-foreground">
                        Change Role
                      </p>
                      {assignableRoles.map((r) => (
                        <button
                          key={r}
                          disabled={r === m.role}
                          onClick={() => {
                            onChangeRole(m.user_id, r);
                            setOpenMenu(null);
                          }}
                          className={cn(
                            'flex w-full items-center gap-2 px-3 py-1.5 text-sm',
                            r === m.role
                              ? 'font-medium text-primary cursor-default'
                              : 'text-foreground hover:bg-accent',
                          )}
                        >
                          <Shield className="h-3 w-3" />
                          {roleLabels[r]}
                          {r === m.role && (
                            <span className="ml-auto text-xs text-muted-foreground">current</span>
                          )}
                        </button>
                      ))}
                      <div className="my-1 border-t border-border" />
                      <button
                        onClick={() => {
                          onRemove(m.user_id);
                          setOpenMenu(null);
                        }}
                        className="flex w-full items-center gap-2 px-3 py-2 text-sm text-destructive hover:bg-destructive/10"
                      >
                        <UserMinus className="h-3 w-3" /> Remove
                      </button>
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
