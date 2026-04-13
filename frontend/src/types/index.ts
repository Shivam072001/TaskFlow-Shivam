export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: { id: string; name: string; email: string };
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  created_by: string;
  created_at: string;
}

export type OrgRole = 'owner' | 'admin' | 'manager' | 'member';

export interface OrgMember {
  id: string;
  org_id: string;
  user_id: string;
  role: OrgRole;
  joined_at: string;
  user_name: string;
  user_email: string;
}

export interface Workspace {
  id: string;
  org_id: string;
  name: string;
  description: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export type WorkspaceRole = 'owner' | 'admin' | 'manager' | 'lead' | 'member' | 'guest' | 'viewer';

export interface WorkspaceMember {
  id: string;
  workspace_id: string;
  user_id: string;
  role: WorkspaceRole;
  joined_at: string;
  user_name: string;
  user_email: string;
}

export interface WorkspaceInvitation {
  id: string;
  workspace_id: string;
  inviter_id: string;
  invitee_email: string;
  invitee_id: string;
  role: WorkspaceRole;
  status: 'pending' | 'accepted' | 'declined';
  created_at: string;
  responded_at: string;
  workspace_name: string;
  inviter_name: string;
}

export interface WorkspaceStats {
  project_count: number;
  tasks_by_status: Record<string, number>;
  tasks_by_assignee: { user_id: string; user_name: string; count: number }[];
  overdue_count: number;
}

export interface Project {
  id: string;
  org_id: string;
  name: string;
  prefix: string;
  description: string;
  workspace_id: string;
  owner_id: string;
  created_at: string;
}

export type TaskStatus = 'todo' | 'in_progress' | 'done' | 'blocked';
export type TaskPriority = 'low' | 'medium' | 'high';

export interface Task {
  id: string;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  project_id: string;
  assignee_id: string | null;
  start_date: string | null;
  due_date: string | null;
  created_by: string;
  created_at: string;
  updated_at: string;
  task_number: number;
  task_key: string;
  blocked_reason: string;
  blocked_by_task: string;
  custom_fields?: Record<string, string>;
}

export interface WIPLimit {
  id: string;
  project_id: string;
  status: TaskStatus;
  max_tasks: number;
}

export interface CustomFieldDefinition {
  id: string;
  project_id: string;
  name: string;
  field_type: 'text' | 'number' | 'select';
  options: string[];
  required: boolean;
  created_by: string;
  created_at: string;
}

export interface CustomFieldValue {
  id: string;
  task_id: string;
  field_id: string;
  value: string;
}

export interface Comment {
  id: string;
  entity_type: 'project' | 'task';
  entity_id: string;
  user_id: string;
  parent_id: string | null;
  content: string;
  created_at: string;
  updated_at: string;
  user_name: string;
  user_email: string;
  replies: Comment[];
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: PaginationMeta;
}

export interface ApiError {
  error: string;
  fields?: Record<string, string>;
}

export type ProjectRole = WorkspaceRole;

export interface Team {
  id: string;
  org_id: string;
  name: string;
  created_by: string;
  created_at: string;
}

export interface TeamMember {
  id: string;
  team_id: string;
  user_id: string;
  added_at: string;
  user_name: string;
  user_email: string;
}

export interface OrgInvitation {
  id: string;
  org_id: string;
  inviter_id: string;
  invitee_email: string;
  invitee_id: string | null;
  role: OrgRole;
  status: 'pending' | 'accepted' | 'declined';
  created_at: string;
  responded_at: string | null;
  org_name: string;
  inviter_name: string;
}

export interface ProjectMember {
  id: string;
  project_id: string;
  user_id: string;
  role: ProjectRole;
  joined_at: string;
  user_name: string;
  user_email: string;
}
