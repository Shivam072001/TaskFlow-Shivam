import { useState, useCallback, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import { Plus, Loader2, Kanban, Mail, Search, X } from 'lucide-react';
import { PageShell } from '@components/layout/PageShell';
import { Breadcrumbs } from '@components/layout/Breadcrumbs';
import { StatsBar } from '@features/workspaces/StatsBar';
import { MemberList } from '@features/workspaces/MemberList';
import { ProjectCard } from '@features/projects/ProjectCard';
import { CreateProjectDialog } from '@features/projects/CreateProjectDialog';
import { Pagination } from '@components/ui/Pagination';
import { useWorkspace, useWorkspaceStats, useWorkspaceMembers, useUpdateMemberRole, useRemoveMember, useDirectAddMember, useLeaveWorkspace } from '@hooks/useWorkspaces';
import { useProjects, useCreateProject } from '@hooks/useProjects';
import { useOrgPrefixes, useOrgMembers } from '@hooks/useOrganizations';
import { useTeams, useAddTeamToWorkspace } from '@hooks/useTeams';
import { useSendInvite, useWorkspaceInvitations } from '@hooks/useInvitations';
import { useAuth } from '@hooks/useAuth';
import type { WorkspaceRole } from '@/types';
import { cn } from '@utils/cn';
import * as wipApi from '@core/api/customFields';

export function WorkspaceDashboardPage() {
  const { orgSlug, wid } = useParams<{ orgSlug: string; wid: string }>();
  const workspaceId = wid!;
  const { activeOrg, user } = useAuth();
  const orgId = activeOrg?.id ?? '';

  const [showCreateProject, setShowCreateProject] = useState(false);

  const [projPage, setProjPage] = useState(1);
  const [projPageSize, setProjPageSize] = useState(5);
  const [searchInput, setSearchInput] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [ownerFilter, setOwnerFilter] = useState('');

  const projectParams = useMemo(() => ({
    page: projPage,
    limit: projPageSize,
    ...(searchQuery && { search: searchQuery }),
    ...(ownerFilter && { owner: ownerFilter }),
  }), [projPage, projPageSize, searchQuery, ownerFilter]);

  const { data: workspace, isLoading: wsLoading, isError: wsError } = useWorkspace(orgId, workspaceId);
  const { data: stats } = useWorkspaceStats(orgId, workspaceId);
  const { data: members = [] } = useWorkspaceMembers(orgId, workspaceId);
  const { data: projData, isLoading: projLoading, isError: projError } = useProjects(orgId, workspaceId, projectParams);
  const { data: pendingInvitations = [] } = useWorkspaceInvitations(orgId, workspaceId);

  const projects = projData?.data ?? [];
  const projTotal = projData?.meta.total ?? 0;

  const { data: existingPrefixes = [] } = useOrgPrefixes(orgId);

  const { data: orgMembersList } = useOrgMembers(orgId);
  const { data: orgTeams } = useTeams(orgId);

  const createProject = useCreateProject(orgId, workspaceId);
  const sendInvite = useSendInvite(orgId, workspaceId);
  const updateRole = useUpdateMemberRole(orgId, workspaceId);
  const removeMember = useRemoveMember(orgId, workspaceId);
  const directAdd = useDirectAddMember(orgId, workspaceId);
  const addTeamToWs = useAddTeamToWorkspace(orgId, workspaceId);
  const leaveWs = useLeaveWorkspace(orgId, workspaceId);

  const currentMember = members.find((m) => m.user_id === user?.id);
  const currentRole = (currentMember?.role ?? 'member') as WorkspaceRole;

  const pendingCount = pendingInvitations.filter((i) => i.status === 'pending').length;

  const memberOptions = useMemo(
    () => [
      { value: '', label: 'All members' },
      ...members.map((m) => ({ value: m.user_id, label: m.user_name })),
    ],
    [members],
  );

  const handleSearch = useCallback(() => {
    setSearchQuery(searchInput.trim());
    setProjPage(1);
  }, [searchInput]);

  const handleClearSearch = useCallback(() => {
    setSearchInput('');
    setSearchQuery('');
    setProjPage(1);
  }, []);

  const handleDirectAdd = useCallback(
    (userId: string, role: string) => directAdd.mutate({ userId, role }),
    [directAdd],
  );

  const handleAddTeam = useCallback(
    (teamId: string, defaultRole: string) => addTeamToWs.mutate({ teamId, defaultRole }),
    [addTeamToWs],
  );

  const handleLeaveWorkspace = useCallback(
    () => leaveWs.mutate(undefined, { onSuccess: () => window.location.href = `/org/${orgSlug}/workspaces` }),
    [leaveWs, orgSlug],
  );

  const handleInvite = useCallback(
    (email: string, role: string) => sendInvite.mutate({ email, role: role as WorkspaceRole }),
    [sendInvite],
  );

  const handleChangeRole = useCallback(
    (userId: string, role: string) => updateRole.mutate({ userId, role }),
    [updateRole],
  );

  const handleRemove = useCallback(
    (userId: string) => removeMember.mutate(userId),
    [removeMember],
  );

  if (wsLoading) {
    return (
      <PageShell>
        <div className="flex justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      </PageShell>
    );
  }

  if (wsError) {
    return (
      <PageShell>
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-6 text-center">
          <p className="text-destructive font-medium">Failed to load workspace</p>
          <p className="mt-1 text-sm text-muted-foreground">The workspace may not exist or you may not have access.</p>
        </div>
      </PageShell>
    );
  }

  return (
    <PageShell>
      <Breadcrumbs items={[
        { label: 'Workspaces', href: `/org/${orgSlug}/workspaces` },
        { label: workspace?.name ?? '' },
      ]} />

      <div className="mb-6">
        <h1 className="text-2xl font-bold text-foreground">{workspace?.name}</h1>
        {workspace?.description && <p className="mt-1 text-muted-foreground">{workspace.description}</p>}
      </div>

      {stats && <div className="mb-8"><StatsBar stats={stats} /></div>}

      <div className="grid gap-8 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-foreground">Projects</h2>
            <button
              onClick={() => setShowCreateProject(true)}
              className={cn('inline-flex items-center gap-1 rounded-lg bg-primary/10 px-3 py-1.5 text-xs font-medium text-primary hover:bg-primary/20')}
            >
              <Plus className="h-3 w-3" /> New Project
            </button>
          </div>

          <div className="mb-4 flex flex-col gap-3 sm:flex-row">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type="text"
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="Search projects..."
                className="w-full rounded-lg border border-border bg-card py-2 pl-9 pr-9 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              />
              {searchInput && (
                <button
                  onClick={handleClearSearch}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                >
                  <X className="h-4 w-4" />
                </button>
              )}
            </div>

            <select
              value={ownerFilter}
              onChange={(e) => {
                setOwnerFilter(e.target.value);
                setProjPage(1);
              }}
              className="rounded-lg border border-border bg-card px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            >
              {memberOptions.map((opt) => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </select>
          </div>

          {projLoading && (
            <div className="flex justify-center py-10">
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
            </div>
          )}

          {!projLoading && projError && (
            <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-6 text-center">
              <p className="text-destructive text-sm">Failed to load projects. Please try again.</p>
            </div>
          )}

          {!projLoading && !projError && projects.length === 0 && (
            <div className="rounded-xl border border-dashed border-border p-8 text-center">
              <Kanban className="mx-auto h-10 w-10 text-muted-foreground/50" />
              <h3 className="mt-3 font-medium text-foreground">
                {searchQuery || ownerFilter ? 'No matching projects' : 'No projects yet'}
              </h3>
              <p className="mt-1 text-sm text-muted-foreground">
                {searchQuery || ownerFilter
                  ? 'Try adjusting your search or filter'
                  : 'Create a project to start managing tasks'}
              </p>
            </div>
          )}

          {projects.length > 0 && (
            <>
              <div className="grid gap-3 sm:grid-cols-2">
                {projects.map((p) => (
                  <ProjectCard key={p.id} project={p} workspaceId={workspaceId} orgSlug={orgSlug!} />
                ))}
              </div>

              {projTotal > projPageSize && (
                <div className="mt-4">
                  <Pagination
                    page={projPage}
                    pageSize={projPageSize}
                    total={projTotal}
                    onPageChange={setProjPage}
                    onPageSizeChange={setProjPageSize}
                  />
                </div>
              )}
            </>
          )}
        </div>

        <div className="space-y-6">
          <MemberList
            members={members}
            orgMembers={orgMembersList}
            orgTeams={orgTeams}
            currentUserRole={currentRole}
            onDirectAdd={handleDirectAdd}
            onAddTeam={handleAddTeam}
            onInvite={handleInvite}
            onChangeRole={handleChangeRole}
            onRemove={handleRemove}
            onLeave={handleLeaveWorkspace}
            isInviting={sendInvite.isPending}
          />

          {pendingCount > 0 && (
            <div className="rounded-xl border border-border bg-card p-4">
              <h3 className="mb-3 flex items-center gap-2 text-sm font-semibold text-foreground">
                <Mail className="h-4 w-4" />
                Pending Invitations
                <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs font-normal text-primary">
                  {pendingCount}
                </span>
              </h3>
              <div className="space-y-2">
                {pendingInvitations
                  .filter((i) => i.status === 'pending')
                  .map((inv) => (
                    <div key={inv.id} className="flex items-center justify-between rounded-md border border-border bg-muted/50 px-3 py-2 text-sm">
                      <span className="truncate text-foreground">{inv.invitee_email}</span>
                      <span className="ml-2 shrink-0 rounded-full bg-yellow-100 px-2 py-0.5 text-xs font-medium text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400">
                        Pending
                      </span>
                    </div>
                  ))}
              </div>
            </div>
          )}
        </div>
      </div>

      <CreateProjectDialog
        open={showCreateProject}
        existingPrefixes={existingPrefixes}
        onClose={() => setShowCreateProject(false)}
        onSubmit={async ({ name, prefix, description, wipLimits }) => {
          createProject.mutate(
            { name, prefix, description },
            {
              onSuccess: async (project) => {
                for (const wl of wipLimits) {
                  await wipApi.setWIPLimit(project.id, wl.status, wl.maxTasks);
                }
                setShowCreateProject(false);
              },
            },
          );
        }}
        isLoading={createProject.isPending}
      />
    </PageShell>
  );
}
