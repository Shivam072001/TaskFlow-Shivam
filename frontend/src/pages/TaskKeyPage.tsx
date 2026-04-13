import { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import { PageShell } from '@components/layout/PageShell';
import { useAuth } from '@hooks/useAuth';
import { useTaskByKey } from '@hooks/useOrganizations';

export function TaskKeyPage() {
  const { orgSlug, taskKey } = useParams<{ orgSlug: string; taskKey: string }>();
  const navigate = useNavigate();
  const { activeOrg } = useAuth();

  const orgId = activeOrg?.slug === orgSlug ? activeOrg.id : '';
  const { data: taskData, isLoading, error } = useTaskByKey(orgId, taskKey ?? '');

  useEffect(() => {
    if (taskData) {
      navigate(
        `/org/${orgSlug}/workspaces/${taskData.workspace_id}/projects/${taskData.project_id}?task=${taskData.task_key}`,
        { replace: true },
      );
    }
  }, [taskData, navigate, orgSlug]);

  if (!orgId) {
    return (
      <PageShell>
        <div className="flex flex-col items-center justify-center py-20">
          <p className="text-muted-foreground">Organization not found. Please select an organization first.</p>
          <button
            onClick={() => navigate('/organizations')}
            className="mt-4 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            Go to Organizations
          </button>
        </div>
      </PageShell>
    );
  }

  if (isLoading) {
    return (
      <PageShell>
        <div className="flex justify-center py-20">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      </PageShell>
    );
  }

  if (error) {
    return (
      <PageShell>
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-6 text-center">
          <p className="text-destructive">Task not found</p>
          <button
            onClick={() => navigate(`/org/${orgSlug}/workspaces`)}
            className="mt-4 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
          >
            Back to Workspaces
          </button>
        </div>
      </PageShell>
    );
  }

  return null;
}
