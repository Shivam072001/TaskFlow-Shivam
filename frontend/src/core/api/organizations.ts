import type { Organization, OrgMember, Task } from '@/types';
import apiClient from './client';

export async function listOrganizations(): Promise<Organization[]> {
  const { data } = await apiClient.get<{ organizations: Organization[] }>('/organizations');
  return data.organizations;
}

export async function createOrganization(name: string, slug: string): Promise<Organization> {
  const { data } = await apiClient.post<Organization>('/organizations', { name, slug });
  return data;
}

export async function getOrganization(orgId: string): Promise<Organization> {
  const { data } = await apiClient.get<Organization>(`/organizations/${orgId}`);
  return data;
}

export async function updateOrganization(orgId: string, body: { name?: string }): Promise<Organization> {
  const { data } = await apiClient.patch<Organization>(`/organizations/${orgId}`, body);
  return data;
}

export async function listOrgMembers(orgId: string): Promise<OrgMember[]> {
  const { data } = await apiClient.get<{ members: OrgMember[] }>(`/organizations/${orgId}/members`);
  return data.members;
}

export async function listOrgPrefixes(orgId: string): Promise<string[]> {
  const { data } = await apiClient.get<{ prefixes: string[] }>(`/organizations/${orgId}/prefixes`);
  return data.prefixes;
}

export interface TaskWithContext extends Task {
  workspace_id: string;
}

export async function getTaskByKey(orgId: string, taskKey: string): Promise<TaskWithContext> {
  const { data } = await apiClient.get<TaskWithContext>(`/organizations/${orgId}/tasks/by-key/${taskKey}`);
  return data;
}
