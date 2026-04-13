import { useState, useCallback, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import { Plus, FolderOpen, Loader2 } from 'lucide-react';
import { PageShell } from '@components/layout/PageShell';
import { WorkspaceCard } from '@features/workspaces/WorkspaceCard';
import { CreateWorkspaceDialog } from '@features/workspaces/CreateWorkspaceDialog';
import { InvitationBanner } from '@features/workspaces/InvitationBanner';
import { Pagination } from '@components/ui/Pagination';
import { useWorkspaces, useCreateWorkspace } from '@hooks/useWorkspaces';
import { useMyInvitations, useRespondToInvitation } from '@hooks/useInvitations';
import { useOrganizations } from '@hooks/useOrganizations';
import { useAuth } from '@hooks/useAuth';
import { cn } from '@utils/cn';

export function WorkspaceListPage() {
  const { orgSlug } = useParams<{ orgSlug: string }>();
  const navigate = useNavigate();
  const { activeOrg, setActiveOrg } = useAuth();
  const { data: orgs } = useOrganizations();

  useEffect(() => {
    if (orgs && orgSlug && (!activeOrg || activeOrg.slug !== orgSlug)) {
      const org = orgs.find((o) => o.slug === orgSlug);
      if (org) setActiveOrg(org);
      else navigate('/organizations', { replace: true });
    }
  }, [orgs, orgSlug, activeOrg, setActiveOrg, navigate]);

  const orgId = activeOrg?.slug === orgSlug ? activeOrg.id : '';

  const [showCreate, setShowCreate] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(12);
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useWorkspaces(orgId, { page, limit: pageSize });
  const createMutation = useCreateWorkspace(orgId);
  const { data: invitations = [] } = useMyInvitations();
  const respondMutation = useRespondToInvitation();

  const workspaces = data?.data ?? [];
  const total = data?.meta.total ?? 0;

  const handleCreate = useCallback(
    (name: string, description: string) => {
      createMutation.mutate({ name, description }, { onSuccess: () => setShowCreate(false) });
    },
    [createMutation],
  );

  const handleInvitationRespond = useCallback(
    (invitationId: string, accept: boolean) => {
      respondMutation.mutate(
        { invitationId, accept },
        { onSuccess: () => { if (accept) queryClient.invalidateQueries({ queryKey: ['workspaces'] }); } },
      );
    },
    [respondMutation, queryClient],
  );

  return (
    <PageShell>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Workspaces</h1>
          <p className="text-muted-foreground text-sm mt-1">Select a workspace or create a new one</p>
        </div>
        <button
          onClick={() => setShowCreate(true)}
          className={cn(
            'inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground',
            'hover:bg-primary/90 transition-colors'
          )}
        >
          <Plus className="h-4 w-4" /> New Workspace
        </button>
      </div>

      <InvitationBanner
        invitations={invitations}
        onRespond={handleInvitationRespond}
        isResponding={respondMutation.isPending}
      />

      {invitations.filter((i) => i.status === 'pending').length > 0 && <div className="mb-6" />}

      {isLoading && (
        <div className="flex justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {error && (
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-6 text-center">
          <p className="text-destructive">Failed to load workspaces</p>
        </div>
      )}

      {!isLoading && workspaces.length === 0 && (
        <div className="rounded-xl border border-dashed border-border p-12 text-center">
          <FolderOpen className="mx-auto h-12 w-12 text-muted-foreground/50" />
          <h3 className="mt-4 font-medium text-foreground">No workspaces yet</h3>
          <p className="mt-1 text-sm text-muted-foreground">Create your first workspace to get started</p>
        </div>
      )}

      {workspaces.length > 0 && (
        <>
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {workspaces.map((ws) => (
              <WorkspaceCard key={ws.id} workspace={ws} orgSlug={orgSlug!} />
            ))}
          </div>

          {total > pageSize && (
            <div className="mt-6">
              <Pagination
                page={page}
                pageSize={pageSize}
                total={total}
                onPageChange={setPage}
                onPageSizeChange={setPageSize}
              />
            </div>
          )}
        </>
      )}

      <CreateWorkspaceDialog
        open={showCreate}
        onClose={() => setShowCreate(false)}
        onSubmit={handleCreate}
        isLoading={createMutation.isPending}
      />
    </PageShell>
  );
}
