# TaskFlow Database Schema

**Database:** PostgreSQL 16  
**Extension:** `pgcrypto` (for `gen_random_uuid()`)  
**Migration tool:** golang-migrate (14 versioned migration pairs)

---

## Table of Contents

- [Entity Relationship Diagram](#entity-relationship-diagram)
- [Tables](#tables)
  - [users](#users)
  - [organizations](#organizations)
  - [organization\_members](#organization_members)
  - [workspaces](#workspaces)
  - [workspace\_members](#workspace_members)
  - [workspace\_invitations](#workspace_invitations)
  - [projects](#projects)
  - [project\_members](#project_members)
  - [tasks](#tasks)
  - [comments](#comments)
  - [custom\_field\_definitions](#custom_field_definitions)
  - [custom\_field\_values](#custom_field_values)
  - [project\_wip\_limits](#project_wip_limits)
  - [teams](#teams)
  - [team\_members](#team_members)
  - [org\_invitations](#org_invitations)
  - [workspace\_teams](#workspace_teams)
  - [project\_teams](#project_teams)
- [Indexes](#indexes)
- [Constraints Summary](#constraints-summary)
- [Role Enumerations](#role-enumerations)

---

## Entity Relationship Diagram

```
                                ┌──────────────┐
                                │    users     │
                                │──────────────│
                                │ id       (PK)│
                                │ name         │
                                │ email   (UQ) │
                                │ password     │
                                │ created_at   │
                                └──────┬───────┘
                                       │
             ┌─────────────────────────┼──────────────────────────┐
             │                         │                          │
             ▼                         ▼                          ▼
   ┌─────────────────┐     ┌───────────────────┐      ┌────────────────────┐
   │  organizations  │     │  org_invitations  │      │       teams        │
   │─────────────────│     │───────────────────│      │────────────────────│
   │ id          (PK)│     │ id           (PK) │      │ id            (PK) │
   │ name            │     │ org_id        (FK)│──┐   │ org_id         (FK)│──┐
   │ slug       (UQ) │     │ inviter_id    (FK)│  │   │ name               │  │
   │ created_by (FK) │     │ invitee_email     │  │   │ created_by     (FK)│  │
   │ created_at      │     │ invitee_id    (FK)│  │   │ created_at         │  │
   └────────┬────────┘     │ role              │  │   └─────────┬──────────┘  │
            │              │ status            │  │             │             │
            │              │ created_at        │  │             ▼             │
            │              │ responded_at      │  │   ┌────────────────────┐  │
            │              └───────────────────┘  │   │   team_members     │  │
            │                                     │   │────────────────────│  │
            ▼                                     │   │ id            (PK) │  │
  ┌──────────────────────┐                        │   │ team_id        (FK)│  │
  │ organization_members │                        │   │ user_id        (FK)│  │
  │──────────────────────│                        │   │ added_at           │  │
  │ id              (PK) │                        │   └────────────────────┘  │
  │ org_id          (FK) │────────────────────────┘                           │
  │ user_id         (FK) │                                                    │
  │ role                 │                                                    │
  │ joined_at            │                                                    │
  └──────────────────────┘                                                    │
            │                                                                 │
            ▼                                                                 │
   ┌─────────────────┐                    ┌─────────────────────┐             │
   │   workspaces    │                    │  workspace_teams    │             │
   │─────────────────│                    │─────────────────────│             │
   │ id         (PK) │◄────────────────── │ workspace_id  (FK)  │             │
   │ name            │                    │ team_id        (FK) │─────────────┘
   │ description     │                    │ default_role        │
   │ org_id     (FK) │                    │ added_at            │
   │ created_by (FK) │                    └─────────────────────┘
   │ created_at      │
   │ updated_at      │
   └────────┬────────┘
            │
    ┌───────┼──────────────────────────────────┐
    │       │                                  │
    ▼       ▼                                  ▼
┌────────────────────┐              ┌──────────────────────┐
│ workspace_members  │              │ workspace_invitations│
│────────────────────│              │──────────────────────│
│ id            (PK) │              │ id              (PK) │
│ workspace_id  (FK) │              │ workspace_id    (FK) │
│ user_id       (FK) │              │ inviter_id      (FK) │
│ role               │              │ invitee_email        │
│ joined_at          │              │ invitee_id      (FK) │
└────────────────────┘              │ role                 │
    │                               │ status               │
    ▼                               │ created_at           │
┌─────────────────┐                 │ responded_at         │
│    projects     │                 └──────────────────────┘
│─────────────────│
│ id         (PK) │
│ name            │         ┌────────────────────┐
│ description     │         │  project_members   │
│ workspace_id(FK)│         │────────────────────│
│ org_id     (FK) │         │ id            (PK) │
│ owner_id   (FK) │         │ project_id    (FK) │──┐
│ prefix          │         │ user_id       (FK) │  │
│ created_at      │         │ role               │  │
└────────┬────────┘         │ joined_at          │  │
         │                  └────────────────────┘  │
         │                                          │
         │      ┌────────────────────┐              │
         │      │   project_teams    │              │
         │      │────────────────────│              │
         ├─────►│ project_id    (FK) │◄─────────────┘
         │      │ team_id       (FK) │
         │      │ default_role       │
         │      │ added_at           │
         │      └────────────────────┘
         │
    ┌────┼───────────────────────────────────┐
    │    │                                   │
    ▼    ▼                                   ▼
┌─────────────────┐         ┌──────────────────────────┐
│     tasks       │         │   project_wip_limits     │
│─────────────────│         │──────────────────────────│
│ id         (PK) │         │ id            (PK)       │
│ title           │         │ project_id    (FK)       │
│ description     │         │ status                   │
│ status          │         │ max_tasks                │
│ priority        │         └──────────────────────────┘
│ project_id (FK) │
│ assignee_id(FK) │         ┌──────────────────────────┐
│ start_date      │         │ custom_field_definitions │
│ due_date        │         │──────────────────────────│
│ created_by (FK) │         │ id            (PK)       │
│ created_at      │         │ project_id    (FK)       │
│ updated_at      │         │ name                     │
│ task_number     │         │ field_type               │
│ task_key        │         │ options        (JSONB)   │
│ blocked_reason  │         │ required                 │
│ blocked_by_task │         │ created_by    (FK)       │
└────────┬────────┘         │ created_at               │
         │                  └──────────────────────────┘
    ┌────┼────────────────┐
    │    │                │
    ▼    ▼                ▼
┌────────────────┐  ┌──────────────────────┐
│   comments     │  │ custom_field_values  │
│────────────────│  │──────────────────────│
│ id        (PK) │  │ id        (PK)       │
│ entity_type    │  │ task_id   (FK)       │
│ entity_id      │  │ field_id  (FK)       │
│ user_id   (FK) │  │ value                │
│ parent_id (FK) │  └──────────────────────┘
│ content        │
│ created_at     │
│ updated_at     │
└────────────────┘
```

---

## Tables

### users

Primary identity table for all authenticated accounts.

| Column       | Type         | Nullable | Default              | Description            |
|--------------|--------------|----------|----------------------|------------------------|
| `id`         | `UUID`       | NO       | `gen_random_uuid()`  | Primary key            |
| `name`       | `TEXT`       | NO       |                      | Display name           |
| `email`      | `TEXT`       | NO       |                      | Login email (unique)   |
| `password`   | `TEXT`       | NO       |                      | bcrypt hash            |
| `created_at` | `TIMESTAMPTZ`| NO       | `now()`              | Registration timestamp |

**Unique:** `email`

---

### organizations

Top-level tenant entity. Every workspace belongs to exactly one organization.

| Column       | Type         | Nullable | Default              | Description                   |
|--------------|--------------|----------|----------------------|-------------------------------|
| `id`         | `UUID`       | NO       | `gen_random_uuid()`  | Primary key                   |
| `name`       | `TEXT`       | NO       |                      | Organization display name     |
| `slug`       | `TEXT`       | NO       |                      | URL-safe identifier (unique)  |
| `created_by` | `UUID`       | NO       |                      | FK → `users.id`               |
| `created_at` | `TIMESTAMPTZ`| NO       | `now()`              | Creation timestamp            |

**Unique:** `slug`

---

### organization_members

Maps users to organizations with a role.

| Column     | Type         | Nullable | Default              | Description           |
|-----------|-------------|----------|----------------------|-----------------------|
| `id`      | `UUID`      | NO       | `gen_random_uuid()`  | Primary key           |
| `org_id`  | `UUID`      | NO       |                      | FK → `organizations.id` (CASCADE) |
| `user_id` | `UUID`      | NO       |                      | FK → `users.id` (CASCADE) |
| `role`    | `TEXT`      | NO       |                      | `owner` · `admin` · `manager` · `member` |
| `joined_at`| `TIMESTAMPTZ`| NO     | `now()`              | Membership timestamp  |

**Unique:** `(org_id, user_id)`

---

### workspaces

A workspace is a container for projects within an organization.

| Column       | Type         | Nullable | Default              | Description                    |
|-------------|-------------|----------|----------------------|--------------------------------|
| `id`        | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                    |
| `name`      | `TEXT`      | NO       |                      | Workspace name                 |
| `description`| `TEXT`     | NO       | `''`                 | Optional description           |
| `org_id`    | `UUID`      | NO       |                      | FK → `organizations.id` (CASCADE) |
| `created_by`| `UUID`      | NO       |                      | FK → `users.id` (CASCADE)      |
| `created_at`| `TIMESTAMPTZ`| NO      | `now()`              | Creation timestamp             |
| `updated_at`| `TIMESTAMPTZ`| NO      | `now()`              | Last update timestamp          |

---

### workspace_members

Maps users to workspaces with a granular role.

| Column         | Type         | Nullable | Default              | Description                        |
|---------------|-------------|----------|----------------------|------------------------------------|
| `id`          | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                        |
| `workspace_id`| `UUID`      | NO       |                      | FK → `workspaces.id` (CASCADE)     |
| `user_id`     | `UUID`      | NO       |                      | FK → `users.id` (CASCADE)          |
| `role`        | `TEXT`      | NO       |                      | `owner` · `admin` · `manager` · `lead` · `member` · `guest` · `viewer` |
| `joined_at`   | `TIMESTAMPTZ`| NO      | `now()`              | Membership timestamp               |

**Unique:** `(workspace_id, user_id)`

---

### workspace_invitations

Pending/resolved invitations to join a workspace.

| Column         | Type         | Nullable | Default              | Description                                  |
|---------------|-------------|----------|----------------------|----------------------------------------------|
| `id`          | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                                  |
| `workspace_id`| `UUID`      | NO       |                      | FK → `workspaces.id` (CASCADE)               |
| `inviter_id`  | `UUID`      | NO       |                      | FK → `users.id` (CASCADE)                    |
| `invitee_email`| `TEXT`     | NO       |                      | Email of the invited user                    |
| `invitee_id`  | `UUID`      | YES      |                      | FK → `users.id` (CASCADE), set on accept     |
| `role`        | `TEXT`      | NO       |                      | `admin` · `manager` · `lead` · `member` · `guest` · `viewer` |
| `status`      | `TEXT`      | NO       | `'pending'`          | `pending` · `accepted` · `declined`          |
| `created_at`  | `TIMESTAMPTZ`| NO      | `now()`              | Invitation timestamp                         |
| `responded_at`| `TIMESTAMPTZ`| YES     |                      | Set when the invitee responds                |

**Unique (partial):** `(workspace_id, invitee_email) WHERE status = 'pending'`

---

### projects

A project lives inside a workspace and holds tasks.

| Column         | Type         | Nullable | Default              | Description                        |
|---------------|-------------|----------|----------------------|------------------------------------|
| `id`          | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                        |
| `name`        | `TEXT`      | NO       |                      | Project name                       |
| `description` | `TEXT`      | NO       | `''`                 | Optional description               |
| `workspace_id`| `UUID`      | NO       |                      | FK → `workspaces.id` (CASCADE)     |
| `org_id`      | `UUID`      | NO       |                      | FK → `organizations.id` (CASCADE)  |
| `owner_id`    | `UUID`      | NO       |                      | FK → `users.id` (CASCADE)          |
| `prefix`      | `TEXT`      | NO       | `''`                 | Task key prefix (e.g. `PROJ`)      |
| `created_at`  | `TIMESTAMPTZ`| NO      | `now()`              | Creation timestamp                 |

**Unique:** `(org_id, prefix)` — prefix is unique within an organization

---

### project_members

Explicit per-project role assignments (independent from workspace roles).

| Column       | Type         | Nullable | Default              | Description                    |
|-------------|-------------|----------|----------------------|--------------------------------|
| `id`        | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                    |
| `project_id`| `UUID`      | NO       |                      | FK → `projects.id` (CASCADE)   |
| `user_id`   | `UUID`      | NO       |                      | FK → `users.id` (CASCADE)      |
| `role`      | `TEXT`      | NO       |                      | `owner` · `admin` · `manager` · `lead` · `member` · `guest` · `viewer` |
| `joined_at` | `TIMESTAMPTZ`| NO      | `now()`              | Membership timestamp           |

**Unique:** `(project_id, user_id)`

---

### tasks

Individual work items within a project.

| Column          | Type         | Nullable | Default              | Description                           |
|----------------|-------------|----------|----------------------|---------------------------------------|
| `id`           | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                           |
| `title`        | `TEXT`      | NO       |                      | Task title                            |
| `description`  | `TEXT`      | NO       | `''`                 | Rich text (HTML) description          |
| `status`       | `TEXT`      | NO       | `'todo'`             | `todo` · `in_progress` · `done` · `blocked` |
| `priority`     | `TEXT`      | NO       | `'medium'`           | `low` · `medium` · `high`             |
| `project_id`   | `UUID`      | NO       |                      | FK → `projects.id` (CASCADE)          |
| `assignee_id`  | `UUID`      | YES      |                      | FK → `users.id` (SET NULL on delete)  |
| `start_date`   | `DATE`      | YES      |                      | Optional start date                   |
| `due_date`     | `DATE`      | YES      |                      | Optional due date                     |
| `created_by`   | `UUID`      | NO       |                      | FK → `users.id` (CASCADE)             |
| `created_at`   | `TIMESTAMPTZ`| NO      | `now()`              | Creation timestamp                    |
| `updated_at`   | `TIMESTAMPTZ`| NO      | `now()`              | Last update timestamp                 |
| `task_number`  | `INT`       | YES      |                      | Auto-incremented per project          |
| `task_key`     | `TEXT`      | YES      |                      | Human-readable key (e.g. `PROJ-42`)   |
| `blocked_reason`| `TEXT`     | NO       | `''`                 | Free-text blocked reason              |
| `blocked_by_task`| `TEXT`    | NO       | `''`                 | Task key that is blocking this task   |

**Unique:** `(project_id, task_number)`

---

### comments

Threaded comments on projects or tasks (polymorphic via `entity_type`).

| Column       | Type         | Nullable | Default              | Description                        |
|-------------|-------------|----------|----------------------|------------------------------------|
| `id`        | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                        |
| `entity_type`| `TEXT`     | NO       |                      | `project` · `task`                 |
| `entity_id` | `UUID`      | NO       |                      | ID of the parent project or task   |
| `user_id`   | `UUID`      | NO       |                      | FK → `users.id` (CASCADE)          |
| `parent_id` | `UUID`      | YES      |                      | FK → `comments.id` (CASCADE), for replies |
| `content`   | `TEXT`      | NO       |                      | Comment body (HTML)                |
| `created_at`| `TIMESTAMPTZ`| NO      | `now()`              | Creation timestamp                 |
| `updated_at`| `TIMESTAMPTZ`| NO      | `now()`              | Last edit timestamp                |

---

### custom_field_definitions

Schema for project-scoped custom fields (text, number, or select).

| Column       | Type         | Nullable | Default              | Description                        |
|-------------|-------------|----------|----------------------|------------------------------------|
| `id`        | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                        |
| `project_id`| `UUID`      | NO       |                      | FK → `projects.id` (CASCADE)       |
| `name`      | `TEXT`      | NO       |                      | Field label                        |
| `field_type`| `TEXT`      | NO       |                      | `text` · `number` · `select`       |
| `options`   | `JSONB`     | YES      | `'[]'`               | Select options array               |
| `required`  | `BOOLEAN`   | NO       | `FALSE`              | Whether the field is mandatory     |
| `created_by`| `UUID`      | NO       |                      | FK → `users.id` (CASCADE)          |
| `created_at`| `TIMESTAMPTZ`| NO      | `now()`              | Creation timestamp                 |

---

### custom_field_values

Actual values for custom fields on individual tasks.

| Column     | Type   | Nullable | Default              | Description                              |
|-----------|--------|----------|----------------------|------------------------------------------|
| `id`      | `UUID` | NO       | `gen_random_uuid()`  | Primary key                              |
| `task_id` | `UUID` | NO       |                      | FK → `tasks.id` (CASCADE)                |
| `field_id`| `UUID` | NO       |                      | FK → `custom_field_definitions.id` (CASCADE) |
| `value`   | `TEXT` | NO       | `''`                 | Stored value (text representation)       |

**Unique:** `(task_id, field_id)`

---

### project_wip_limits

Work-in-progress limits per status column within a project.

| Column       | Type   | Nullable | Default              | Description                    |
|-------------|--------|----------|----------------------|--------------------------------|
| `id`        | `UUID` | NO       | `gen_random_uuid()`  | Primary key                    |
| `project_id`| `UUID` | NO       |                      | FK → `projects.id` (CASCADE)   |
| `status`    | `TEXT` | NO       |                      | `todo` · `in_progress` · `done`|
| `max_tasks` | `INT`  | NO       |                      | Max concurrent tasks (> 0)     |

**Unique:** `(project_id, status)`

---

### teams

Organization-scoped teams that can be assigned to workspaces and projects.

| Column       | Type         | Nullable | Default              | Description                      |
|-------------|-------------|----------|----------------------|----------------------------------|
| `id`        | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                      |
| `org_id`    | `UUID`      | NO       |                      | FK → `organizations.id` (CASCADE)|
| `name`      | `TEXT`      | NO       |                      | Team display name                |
| `created_by`| `UUID`      | NO       |                      | FK → `users.id`                  |
| `created_at`| `TIMESTAMPTZ`| NO      | `now()`              | Creation timestamp               |

**Unique:** `(org_id, name)`

---

### team_members

Maps users to teams.

| Column     | Type         | Nullable | Default              | Description              |
|-----------|-------------|----------|----------------------|--------------------------|
| `id`      | `UUID`      | NO       | `gen_random_uuid()`  | Primary key              |
| `team_id` | `UUID`      | NO       |                      | FK → `teams.id` (CASCADE)|
| `user_id` | `UUID`      | NO       |                      | FK → `users.id` (CASCADE)|
| `added_at`| `TIMESTAMPTZ`| NO      | `now()`              | Addition timestamp       |

**Unique:** `(team_id, user_id)`

---

### org_invitations

Pending/resolved invitations to join an organization.

| Column         | Type         | Nullable | Default              | Description                              |
|---------------|-------------|----------|----------------------|------------------------------------------|
| `id`          | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                              |
| `org_id`      | `UUID`      | NO       |                      | FK → `organizations.id` (CASCADE)        |
| `inviter_id`  | `UUID`      | NO       |                      | FK → `users.id`                          |
| `invitee_email`| `TEXT`     | NO       |                      | Email of the invited user                |
| `invitee_id`  | `UUID`      | YES      |                      | FK → `users.id`, set on accept           |
| `role`        | `TEXT`      | NO       |                      | `admin` · `manager` · `member`           |
| `status`      | `TEXT`      | NO       | `'pending'`          | `pending` · `accepted` · `declined`      |
| `created_at`  | `TIMESTAMPTZ`| NO      | `now()`              | Invitation timestamp                     |
| `responded_at`| `TIMESTAMPTZ`| YES     |                      | Set when the invitee responds            |

**Unique (partial):** `(org_id, invitee_email) WHERE status = 'pending'`

---

### workspace_teams

Junction table — assigns teams to workspaces with a default role for cascaded members.

| Column         | Type         | Nullable | Default              | Description                        |
|---------------|-------------|----------|----------------------|------------------------------------|
| `id`          | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                        |
| `workspace_id`| `UUID`      | NO       |                      | FK → `workspaces.id` (CASCADE)     |
| `team_id`     | `UUID`      | NO       |                      | FK → `teams.id` (CASCADE)          |
| `default_role`| `TEXT`      | NO       | `'member'`           | Role given to cascaded team members|
| `added_at`    | `TIMESTAMPTZ`| NO      | `now()`              | Assignment timestamp               |

**Unique:** `(workspace_id, team_id)`

---

### project_teams

Junction table — assigns teams to projects with a default role for cascaded members.

| Column         | Type         | Nullable | Default              | Description                        |
|---------------|-------------|----------|----------------------|------------------------------------|
| `id`          | `UUID`      | NO       | `gen_random_uuid()`  | Primary key                        |
| `project_id`  | `UUID`      | NO       |                      | FK → `projects.id` (CASCADE)       |
| `team_id`     | `UUID`      | NO       |                      | FK → `teams.id` (CASCADE)          |
| `default_role`| `TEXT`      | NO       | `'member'`           | Role given to cascaded team members|
| `added_at`    | `TIMESTAMPTZ`| NO      | `now()`              | Assignment timestamp               |

**Unique:** `(project_id, team_id)`

---

## Indexes

| Table                      | Index Name                           | Columns / Condition                            |
|---------------------------|--------------------------------------|------------------------------------------------|
| `users`                   | `idx_users_email`                    | `email` (unique)                               |
| `organizations`           | `idx_organizations_slug`             | `slug` (unique)                                |
| `organization_members`    | `idx_org_members_org`                | `org_id`                                       |
| `organization_members`    | `idx_org_members_user`               | `user_id`                                      |
| `workspace_members`       | `idx_workspace_members_workspace`    | `workspace_id`                                 |
| `workspace_members`       | `idx_workspace_members_user`         | `user_id`                                      |
| `workspace_invitations`   | `idx_invitations_workspace`          | `workspace_id`                                 |
| `workspace_invitations`   | `idx_invitations_invitee`            | `invitee_email`                                |
| `workspace_invitations`   | `idx_invitations_invitee_id`         | `invitee_id`                                   |
| `workspace_invitations`   | `idx_invitations_unique_pending`     | `(workspace_id, invitee_email) WHERE status = 'pending'` (unique) |
| `org_invitations`         | `idx_org_inv_unique_pending`         | `(org_id, invitee_email) WHERE status = 'pending'` (unique) |
| `projects`                | `idx_projects_workspace`             | `workspace_id`                                 |
| `projects`                | `idx_projects_owner`                 | `owner_id`                                     |
| `tasks`                   | `idx_tasks_project`                  | `project_id`                                   |
| `tasks`                   | `idx_tasks_assignee`                 | `assignee_id`                                  |
| `tasks`                   | `idx_tasks_status`                   | `status`                                       |
| `comments`                | `idx_comments_entity`                | `(entity_type, entity_id)`                     |
| `comments`                | `idx_comments_parent`                | `parent_id`                                    |
| `comments`                | `idx_comments_user`                  | `user_id`                                      |
| `custom_field_definitions`| `idx_custom_fields_project`          | `project_id`                                   |
| `custom_field_values`     | `idx_custom_field_values_task`       | `task_id`                                      |

---

## Constraints Summary

| Table                    | Constraint                                          | Type       |
|--------------------------|-----------------------------------------------------|------------|
| `organization_members`   | `role IN ('owner','admin','manager','member')`       | CHECK      |
| `workspace_members`      | `role IN ('owner','admin','manager','lead','member','guest','viewer')` | CHECK |
| `workspace_invitations`  | `role IN ('admin','manager','lead','member','guest','viewer')` | CHECK |
| `workspace_invitations`  | `status IN ('pending','accepted','declined')`        | CHECK      |
| `org_invitations`        | `role IN ('admin','manager','member')`                | CHECK      |
| `org_invitations`        | `status IN ('pending','accepted','declined')`         | CHECK      |
| `project_members`        | `role IN ('owner','admin','manager','lead','member','guest','viewer')` | CHECK |
| `tasks`                  | `status IN ('todo','in_progress','done','blocked')`   | CHECK      |
| `tasks`                  | `priority IN ('low','medium','high')`                 | CHECK      |
| `project_wip_limits`     | `status IN ('todo','in_progress','done')`             | CHECK      |
| `project_wip_limits`     | `max_tasks > 0`                                       | CHECK      |
| `custom_field_definitions`| `field_type IN ('text','number','select')`           | CHECK      |
| `comments`               | `entity_type IN ('project','task')`                   | CHECK      |

---

## Role Enumerations

### Organization Roles

```
owner → admin → manager → member
```

| Role      | Permissions                                                  |
|-----------|--------------------------------------------------------------|
| `owner`   | Full control, transfer ownership, delete org                 |
| `admin`   | Manage members, workspaces, teams, invite to org             |
| `manager` | Manage teams, view member stats for roles below              |
| `member`  | Access granted workspaces, basic participation               |

### Workspace Roles

```
owner → admin → manager → lead → member → guest → viewer
```

| Role      | Permissions                                                  |
|-----------|--------------------------------------------------------------|
| `owner`   | Full workspace control, delete workspace                     |
| `admin`   | Manage members, projects, settings, invite members           |
| `manager` | Manage projects, assign tasks                                |
| `lead`    | Lead projects, manage task assignments                       |
| `member`  | Create/update tasks, comment                                 |
| `guest`   | Limited task interaction within assigned projects            |
| `viewer`  | Read-only access                                             |

### Project Roles

Same 7-tier system as workspace roles. Assigned **independently** — a user's project role does not inherit from their workspace role.

```
owner → admin → manager → lead → member → guest → viewer
```
