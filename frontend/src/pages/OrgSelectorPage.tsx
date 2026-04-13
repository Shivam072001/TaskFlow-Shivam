import { useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, Building2, Loader2, Mail, Check, X } from 'lucide-react';
import { PageShell } from '@components/layout/PageShell';
import { useOrganizations, useCreateOrganization } from '@hooks/useOrganizations';
import { useMyOrgInvitations, useRespondOrgInvite } from '@hooks/useOrgInvitations';
import { useAuth } from '@hooks/useAuth';
import type { Organization } from '@/types';
import type { AxiosError } from 'axios';
import { cn } from '@utils/cn';

export function OrgSelectorPage() {
  const navigate = useNavigate();
  const { setActiveOrg } = useAuth();
  const { data: orgs, isLoading, isError: orgsError } = useOrganizations();
  const createOrg = useCreateOrganization();
  const { data: pendingInvites } = useMyOrgInvitations();
  const respondInvite = useRespondOrgInvite();

  const [showCreate, setShowCreate] = useState(false);
  const [name, setName] = useState('');
  const [slug, setSlug] = useState('');
  const [slugError, setSlugError] = useState('');

  const handleSelect = useCallback(
    (org: Organization) => {
      setActiveOrg(org);
      navigate(`/org/${org.slug}/workspaces`);
    },
    [setActiveOrg, navigate],
  );

  const handleCreate = useCallback(() => {
    if (!name.trim() || !slug.trim()) return;
    createOrg.mutate(
      { name: name.trim(), slug: slug.trim().toLowerCase() },
      {
        onSuccess: (org) => {
          setActiveOrg(org);
          navigate(`/org/${org.slug}/workspaces`);
        },
        onError: (err: AxiosError<{ fields?: Record<string, string> }>) => {
          const fields = err?.response?.data?.fields;
          if (fields?.slug) setSlugError(fields.slug);
        },
      },
    );
  }, [name, slug, createOrg, setActiveOrg, navigate]);

  const handleSlugChange = (val: string) => {
    setSlug(val.toLowerCase().replace(/[^a-z0-9-]/g, ''));
    setSlugError('');
  };

  return (
    <PageShell>
      <div className="mx-auto max-w-2xl">
        <div className="text-center mb-8">
          <Building2 className="mx-auto h-12 w-12 text-primary" />
          <h1 className="mt-4 text-2xl font-bold text-foreground">Select an Organization</h1>
          <p className="mt-1 text-muted-foreground">Choose an organization to work in, or create a new one</p>
        </div>

        {pendingInvites && pendingInvites.length > 0 && (
          <div className="mb-6 space-y-3">
            <h2 className="text-sm font-medium text-muted-foreground flex items-center gap-2">
              <Mail className="h-4 w-4" /> Pending Invitations
            </h2>
            {pendingInvites.map(inv => (
              <div key={inv.id} className="flex items-center justify-between rounded-xl border border-primary/20 bg-primary/5 p-4">
                <div>
                  <p className="font-medium text-foreground">{inv.org_name}</p>
                  <p className="text-sm text-muted-foreground">
                    Invited by {inv.inviter_name} as {inv.role}
                  </p>
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => respondInvite.mutate({ invitationId: inv.id, accept: true })}
                    disabled={respondInvite.isPending}
                    className="inline-flex items-center gap-1 rounded-lg bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
                  >
                    <Check className="h-3 w-3" /> Accept
                  </button>
                  <button
                    onClick={() => respondInvite.mutate({ invitationId: inv.id, accept: false })}
                    disabled={respondInvite.isPending}
                    className="inline-flex items-center gap-1 rounded-lg border border-border px-3 py-1.5 text-xs font-medium text-foreground hover:bg-muted transition-colors"
                  >
                    <X className="h-3 w-3" /> Decline
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}

        {isLoading && (
          <div className="flex justify-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        )}

        {!isLoading && orgsError && (
          <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-6 text-center mb-6">
            <p className="text-destructive font-medium">Failed to load organizations</p>
            <p className="mt-1 text-sm text-muted-foreground">Please check your connection and try again.</p>
          </div>
        )}

        {!isLoading && !orgsError && orgs && orgs.length > 0 && (
          <div className="grid gap-3 mb-6">
            {orgs.map((org) => (
              <button
                key={org.id}
                onClick={() => handleSelect(org)}
                className={cn(
                  'flex items-center gap-4 rounded-xl border border-border bg-card p-4 text-left w-full',
                  'hover:border-primary/50 hover:shadow-md transition-all',
                )}
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 shrink-0">
                  <Building2 className="h-5 w-5 text-primary" />
                </div>
                <div className="min-w-0">
                  <h3 className="font-semibold text-foreground">{org.name}</h3>
                  <p className="text-sm text-muted-foreground truncate">{org.slug}</p>
                </div>
              </button>
            ))}
          </div>
        )}

        {!isLoading && (!orgs || orgs.length === 0) && !showCreate && (
          <div className="rounded-xl border border-dashed border-border p-12 text-center mb-6">
            <Building2 className="mx-auto h-10 w-10 text-muted-foreground/50" />
            <h3 className="mt-4 font-medium text-foreground">No organizations yet</h3>
            <p className="mt-1 text-sm text-muted-foreground">Create your first organization to get started</p>
          </div>
        )}

        {!showCreate ? (
          <div className="text-center">
            <button
              onClick={() => setShowCreate(true)}
              className={cn(
                'inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2.5 text-sm font-medium text-primary-foreground',
                'hover:bg-primary/90 transition-colors',
              )}
            >
              <Plus className="h-4 w-4" /> Create Organization
            </button>
          </div>
        ) : (
          <div className="rounded-xl border border-border bg-card p-6">
            <h3 className="text-lg font-semibold text-foreground mb-4">Create Organization</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-foreground mb-1">Name</label>
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="My Company"
                  className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-foreground mb-1">
                  Slug <span className="text-muted-foreground font-normal">(URL identifier)</span>
                </label>
                <input
                  type="text"
                  value={slug}
                  onChange={(e) => handleSlugChange(e.target.value)}
                  placeholder="my-company"
                  className={cn(
                    'w-full rounded-lg border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring',
                    slugError ? 'border-destructive' : 'border-border',
                  )}
                />
                {slugError && <p className="mt-1 text-xs text-destructive">{slugError}</p>}
                <p className="mt-1 text-xs text-muted-foreground">3-30 characters, lowercase letters, numbers, and hyphens</p>
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <button
                  onClick={() => { setShowCreate(false); setName(''); setSlug(''); setSlugError(''); }}
                  className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-foreground hover:bg-muted transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreate}
                  disabled={!name.trim() || !slug.trim() || createOrg.isPending}
                  className={cn(
                    'inline-flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground',
                    'hover:bg-primary/90 transition-colors disabled:opacity-50',
                  )}
                >
                  {createOrg.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
                  Create
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </PageShell>
  );
}
