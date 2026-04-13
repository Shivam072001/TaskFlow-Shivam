import type { Team, TeamMember } from '@/types';
import apiClient from './client';

export async function listTeams(orgId: string): Promise<Team[]> {
  const { data } = await apiClient.get<{ teams: Team[] }>(`/organizations/${orgId}/teams`);
  return data.teams;
}

export async function createTeam(orgId: string, name: string): Promise<Team> {
  const { data } = await apiClient.post<Team>(`/organizations/${orgId}/teams`, { name });
  return data;
}

export async function getTeam(orgId: string, teamId: string): Promise<{ team: Team; members: TeamMember[] }> {
  const { data } = await apiClient.get<{ team: Team; members: TeamMember[] }>(
    `/organizations/${orgId}/teams/${teamId}`,
  );
  return data;
}

export async function updateTeam(orgId: string, teamId: string, name: string): Promise<void> {
  await apiClient.patch(`/organizations/${orgId}/teams/${teamId}`, { name });
}

export async function deleteTeam(orgId: string, teamId: string): Promise<void> {
  await apiClient.delete(`/organizations/${orgId}/teams/${teamId}`);
}

export async function addTeamMember(orgId: string, teamId: string, userId: string): Promise<TeamMember> {
  const { data } = await apiClient.post<TeamMember>(
    `/organizations/${orgId}/teams/${teamId}/members`,
    { user_id: userId },
  );
  return data;
}

export async function removeTeamMember(orgId: string, teamId: string, userId: string): Promise<void> {
  await apiClient.delete(`/organizations/${orgId}/teams/${teamId}/members/${userId}`);
}

export async function addTeamToWorkspace(
  orgId: string, workspaceId: string, teamId: string, defaultRole: string,
): Promise<void> {
  await apiClient.post(
    `/organizations/${orgId}/workspaces/${workspaceId}/teams`,
    { team_id: teamId, default_role: defaultRole },
  );
}

export async function removeTeamFromWorkspace(
  orgId: string, workspaceId: string, teamId: string,
): Promise<void> {
  await apiClient.delete(`/organizations/${orgId}/workspaces/${workspaceId}/teams/${teamId}`);
}

export async function addTeamToProject(
  orgId: string, workspaceId: string, projectId: string, teamId: string, defaultRole: string,
): Promise<void> {
  await apiClient.post(
    `/organizations/${orgId}/workspaces/${workspaceId}/projects/${projectId}/teams`,
    { team_id: teamId, default_role: defaultRole },
  );
}
