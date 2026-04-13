# TaskFlow

A production-grade, multi-tenant task management platform with organizations, workspaces, teams, Jira-style task keys, and a Kanban board with drag-and-drop — built as a full-stack monorepo.

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
  - [Backend: Clean Layered Architecture](#backend-clean-layered-architecture)
  - [Frontend: Feature-Sliced Architecture](#frontend-feature-sliced-architecture)
- [Entity Hierarchy](#entity-hierarchy)
- [Role System](#role-system)
- [Getting Started](#getting-started)
- [Seed Data](#seed-data)
- [Test Credentials](#test-credentials)
- [API Reference](#api-reference)
- [Project Structure](#project-structure)
- [Environment Variables](#environment-variables)
- [What You'd Do With More Time](#what-youd-do-with-more-time)

---

## Overview

TaskFlow lets teams organize work across **Organizations**, **Workspaces**, and **Projects**. Users can be grouped into **Teams**, invited at the organization level, and assigned independent roles at every tier. Tasks live on a Kanban board with four statuses (Todo, In Progress, Blocked, Done), Jira-style keys (`PROJ-001`), custom fields, comments, WIP limits, and full pagination/filtering/search. The app features **real-time updates via WebSocket** (live task changes, automatic session expiry), a **personal dashboard** with cross-project task overview, and an **org-level admin dashboard** where managers can view member stats and tasks. Dark mode, drag-and-drop, and responsive design (375px–1280px) are built in.

---

## Tech Stack

| Layer | Technology |
|---|---|
| **Backend** | Go 1.22, Chi v5, pgx v5, golang-migrate, bcrypt, JWT (HS256), gorilla/websocket |
| **Frontend** | React 19, TypeScript, Vite, Tailwind CSS v4, TanStack Query, Zustand, @dnd-kit, Sonner (toasts) |
| **Component Library** | Custom — built with Tailwind CSS, Lucide React icons, and Sonner for notifications (stated choice: no third-party UI kit) |
| **Database** | PostgreSQL 16 (production-ready with [Neon](https://neon.tech) serverless Postgres) |
| **Infrastructure** | Docker, Docker Compose, Nginx (with WebSocket proxy), multi-stage builds |

---

## Architecture

### Backend: Clean Layered Architecture

The backend follows the **Clean / Layered Architecture** pattern (Handler → Service → Repository), sometimes called **Ports & Adapters** in Go communities. Each layer has a single responsibility and only depends on the layer directly below it — never sideways or upward.

```
HTTP Request
     │
     ▼
┌─────────────────────────────────────────────────┐
│  Middleware (CORS, Logger, Recovery, RealIP)     │
├─────────────────────────────────────────────────┤
│  Auth Middleware (JWT validation)                │
├─────────────────────────────────────────────────┤
│  Guard Middleware (OrgGuard / WorkspaceGuard)    │
└──────────────────────┬──────────────────────────┘
                       │
                       ▼
               ┌──────────────┐
               │   Handler    │  HTTP parsing, validation, response writing
               └──────┬───────┘
                      │
                      ▼
               ┌──────────────┐
               │   Service    │  Business rules, authorization, orchestration
               └──────┬───────┘
                      │
                      ▼
               ┌──────────────┐
               │  Repository  │  Raw SQL via pgx, returns domain models
               └──────┬───────┘
                      │
                      ▼
               ┌──────────────┐
               │  PostgreSQL  │
               └──────────────┘
```

**Layer responsibilities:**

| Layer | Responsibility | Knows About |
|---|---|---|
| **Handler** | Parse HTTP requests, validate input, write JSON responses | Service |
| **Service** | Business rules, authorization checks, cross-repo orchestration | Repository, Domain |
| **Repository** | Raw SQL queries via pgx, map rows to domain structs | Domain, Database |
| **Domain** | Pure data structs and constants — zero dependencies | Nothing |
| **Middleware** | Cross-cutting concerns (auth, guards, logging, recovery) | Service (for auth validation) |

Constructor-based dependency injection wired in `main.go`. No frameworks, no globals, no ORM.

**Why we chose it:**

| Benefit | Detail |
|---|---|
| **Testability** | Each layer can be unit-tested in isolation by swapping the layer below with a mock/stub. Services don't know about HTTP; repositories don't know about business rules. |
| **Separation of concerns** | A handler never touches SQL. A repository never checks permissions. Business logic lives in exactly one place (service), making bugs easy to trace. |
| **Readability** | New contributors can follow a request top-to-bottom: handler → service → repository → SQL. No magic, no auto-wiring. |
| **Flexibility** | Swapping PostgreSQL for another store means rewriting only the repository layer. Switching from Chi to another router means rewriting only handlers. |
| **Go idiom** | This is the dominant pattern in production Go services (used by Google, Uber, Cloudflare). It leverages Go's strengths — explicit code, interfaces, constructor injection — without fighting the language. |

**Tradeoffs and limitations:**

| Tradeoff | Detail |
|---|---|
| **Boilerplate** | Each new domain entity requires a new file in handler, service, repository, domain, and dto. For 10+ entities this adds up, but it keeps each file focused and small. |
| **No interface abstraction** | We use concrete struct types instead of Go interfaces for repositories. This is a deliberate simplicity choice — interfaces would be added when we need multiple implementations or test mocks. |
| **Manual DI wiring** | All dependencies are wired by hand in `main.go`. This is verbose but fully transparent. Frameworks like Wire or Fx could automate this, at the cost of magic. |
| **Cross-service calls** | Some services call other services (e.g., `TeamService` calls `WorkspaceRepo` and `ProjectMemberRepo`). This is pragmatic but can lead to circular dependencies if not carefully managed. |

---

### Frontend: Feature-Sliced Architecture

The frontend follows a **Feature-Sliced** architecture with alias-driven imports, where code is organized by domain feature (tasks, workspaces, projects) rather than technical role. Path aliases (`@pages/`, `@features/`, `@hooks/`, `@core/`) enforce a clear unidirectional dependency flow.

```
  Route (URL)
       │
       ▼
  ┌──────────┐
  │  Pages   │  @pages/ — route-level components (one per URL)
  └────┬─────┘
       │ imports
       ▼
  ┌──────────┐
  │ Features │  @features/ — domain UI (KanbanBoard, MemberList, TaskModal)
  └────┬─────┘
       │ imports
       ▼
  ┌──────────┐
  │  Hooks   │  @hooks/ — TanStack Query hooks (server state management)
  └────┬─────┘
       │ imports
       ▼
  ┌──────────┐
  │ Core API │  @core/api/ — Axios HTTP clients (one file per domain)
  └────┬─────┘
       │
       ▼
  ┌──────────┐
  │  Types   │  @types/ — shared TypeScript interfaces
  └──────────┘
```

**Layer responsibilities:**

| Layer | Responsibility | Examples |
|---|---|---|
| **Pages** | Route-level orchestration — compose features, manage URL params, coordinate data | `WorkspaceDashboardPage`, `ProjectDetailPage`, `OrgSettingsPage` |
| **Features** | Domain-specific UI components with their own state and interactions | `KanbanBoard`, `TaskModal`, `MemberList`, `ProjectMemberList` |
| **Hooks** | TanStack Query wrappers — queries, mutations, cache invalidation | `useTasks`, `useWorkspaces`, `useTeams`, `useProjectMembers` |
| **Core API** | Thin Axios wrappers — one function per API endpoint, typed request/response | `tasks.ts`, `workspaces.ts`, `teams.ts`, `orgInvitations.ts` |
| **Store** | Zustand for client-only state (auth token, current user) | `useAuthStore` |
| **Components/UI** | Reusable, domain-agnostic UI primitives | `Pagination`, `Button`, `Navbar` |

**State management split:**

| State Type | Tool | Why |
|---|---|---|
| **Server state** (tasks, members, projects) | TanStack Query | Automatic caching, background refetching, optimistic updates, deduplication |
| **Client state** (auth token, user) | Zustand | Lightweight, no boilerplate, persisted to localStorage |
| **UI state** (modals, form inputs) | React `useState` | Local and ephemeral — no need for global state |

**Why we chose it:**

| Benefit | Detail |
|---|---|
| **Colocation** | Each feature's UI, logic, and types live close together. Finding "everything about tasks" means looking in `@features/tasks/`, `@hooks/useTasks.ts`, and `@core/api/tasks.ts`. |
| **Unidirectional flow** | Pages → Features → Hooks → API. No circular imports. A hook never imports a page; an API client never imports a component. |
| **Alias-driven boundaries** | Path aliases (`@features/`, `@hooks/`) make imports self-documenting and enforce architectural boundaries at the import level. |
| **Scalability** | Adding a new domain (e.g., "Notifications") means adding one file in each layer — not restructuring the app. The pattern scales linearly with feature count. |
| **Server state done right** | TanStack Query eliminates hand-rolled loading/error/cache logic. Optimistic updates on the Kanban board make drag-and-drop feel instant even on slow networks. |
| **Tiny client store** | Zustand replaces Redux for the small amount of client state (auth only), avoiding massive boilerplate for what amounts to two values (token + user). |

**Tradeoffs and limitations:**

| Tradeoff | Detail |
|---|---|
| **Feature vs. component ambiguity** | Some components blur the line between a "feature" and a reusable "UI component" (e.g., `MemberList` is feature-specific but used in multiple pages). We resolve this by putting domain-coupled components in `@features/` and truly generic ones in `@components/ui/`. |
| **No shared feature-to-feature imports** | Features should not import from each other. When two features need shared logic, it must be lifted to `@hooks/` or `@utils/`. This occasionally forces refactoring. |
| **Hook proliferation** | Each domain gets its own hooks file with 5-10 hooks (queries + mutations). With 10+ domains this becomes many files, but each file is small and focused. |
| **Client-side SPA by design** | TaskFlow is an authenticated dashboard — every route is behind login, so SEO is irrelevant and first-paint is offset by cached bundles on repeat visits. A static Nginx deployment keeps infrastructure simple with zero Node.js runtime overhead, and Vite's code-splitting keeps the initial bundle lean. SSR (Next.js/Remix) would add significant complexity for no measurable user benefit on this type of app. |
| **Lightweight form handling** | Forms use controlled `useState` with explicit field-level validation in DTOs. All forms in TaskFlow are small (2–6 fields) — login, register, create workspace, create project, create task — where `useState` is more readable and has zero added dependency weight compared to React Hook Form (~10KB). The backend enforces the same validation rules server-side, so client forms stay thin by design. |

---

## Entity Hierarchy

```
Organization
├── Members (owner / admin / manager / member)
├── Teams
│   └── Team Members
├── Org Invitations
│
└── Workspace
    ├── Members (owner / admin / manager / lead / member / guest / viewer)
    ├── Workspace Teams
    │
    └── Project
        ├── Members (owner / admin / manager / lead / member / guest / viewer)
        ├── Project Teams
        ├── Custom Fields
        ├── WIP Limits
        ├── Comments
        │
        └── Task
            ├── Jira-style Key (PREFIX-001)
            ├── Status: todo | in_progress | blocked | done
            ├── Priority: low | medium | high
            ├── Assignee
            ├── Blocked reason / Blocked by task
            ├── Start date / Due date
            ├── Comments
            └── Custom Field Values
```

---

## Role System

Roles at each level are **fully independent** — there is no inheritance from parent to child. When adding a user or team to a workspace or project, the role must be explicitly assigned.

### Organization Roles

| Role | Manage Members | Manage Teams | Manage Workspaces | Transfer Ownership |
|---|:---:|:---:|:---:|:---:|
| `owner` | Yes | Yes | Yes | Yes |
| `admin` | Yes | Yes | Yes | No |
| `manager` | No | Yes | Yes | No |
| `member` | No | No | No | No |

### Workspace Roles

| Role | Description |
|---|---|
| `owner` | Full control, can delete workspace |
| `admin` | Manage members and settings |
| `manager` | Manage projects and teams |
| `lead` | Project oversight |
| `member` | Standard contributor |
| `guest` | Limited access |
| `viewer` | Read-only |

### Project Roles

Same role set as Workspace: `owner`, `admin`, `manager`, `lead`, `member`, `guest`, `viewer`.

### Cascading Rules

- A user **must be a workspace member** before they can be added to a project within that workspace.
- **Leaving a workspace** removes the user from all projects in that workspace.
- **Leaving an organization** removes the user from all its workspaces, projects, and teams.

---

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

### Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/your-name/taskflow.git
cd taskflow

# 2. Copy environment file
cp .env.example .env

# 3. Build and start all services
docker compose up
```

### What Happens on Startup

```
docker compose up
         │
         ▼
┌─────────────────────────────────────────────┐
│  1. PostgreSQL starts, waits for healthcheck │
│  2. Backend container starts                 │
│     ├── Runs all 14 database migrations      │
│     ├── Runs Go seeder (10 users, 3 orgs,   │
│     │   3 workspaces, 2 teams, 10 projects,  │
│     │   ~190 tasks)                           │
│     └── Starts the API server on :8080       │
│  3. Frontend container starts (Nginx on :80) │
└─────────────────────────────────────────────┘
```

### Access Points

| Service | URL |
|---|---|
| **Frontend** | http://localhost:3000 |
| **Backend API** | http://localhost:8080 |
| **WebSocket** | `ws://localhost:3000/ws?token=<jwt>` (via Nginx) or `ws://localhost:8080/ws?token=<jwt>` (direct) |
| **PostgreSQL** | `localhost:5432` |
| **Health Check** | http://localhost:8080/health |

### Manual Commands

```bash
# Rebuild and restart
docker compose up

# Stop all services
docker compose down

# Stop and remove all data (clean slate)
docker compose down -v

# View backend logs
docker compose logs -f backend

# Run migrations manually
docker compose exec backend migrate -path /app/migrations -database "$DATABASE_URL" up
```

---

## Seed Data

Seed data is applied **automatically** on every `docker compose up`. The seeder is idempotent (`ON CONFLICT DO NOTHING`) so it is safe to run repeatedly.

### What Gets Created

| Entity | Count | Details |
|---|---|---|
| **Users** | 10 | All with password `password123` |
| **Organizations** | 3 | Acme Corp, Globex Inc, Initech |
| **Workspaces** | 3 | Engineering, Marketing, Operations |
| **Teams** | 2 | Backend Squad, Frontend Squad (in Acme Corp) |
| **Projects** | 10 | 7 in Engineering, 2 in Marketing, 1 in Operations |
| **Tasks** | ~190 | 4-6 per status per project across all 4 statuses |

### Seeded Organizations

| Organization | Workspace | Projects |
|---|---|---|
| **Acme Corp** (`acme-corp`) | Engineering | API Gateway, User Service, Payment Engine, Notification Hub, Admin Dashboard, Mobile App, Data Pipeline |
| **Globex Inc** (`globex-inc`) | Marketing | Marketing Site, CRM Integration |
| **Initech** (`initech`) | Operations | Internal Wiki |

### Seeded Teams (Acme Corp)

| Team | Members |
|---|---|
| **Backend Squad** | Charlie Manager, Diana Lead, Eve Developer |
| **Frontend Squad** | Frank Designer, Grace Tester |

---

## Test Credentials

All 10 users share the same password: **`password123`**

| # | Name | Email | Acme Corp | Globex Inc | Initech | Engineering WS | Marketing WS | Operations WS |
|---|---|---|---|---|---|---|---|---|
| 1 | Alice Owner | `test@example.com` | owner | owner | owner | owner | owner | owner |
| 2 | Bob Admin | `bob@example.com` | admin | — | — | admin | — | — |
| 3 | Charlie Manager | `charlie@example.com` | manager | — | — | manager | — | — |
| 4 | Diana Lead | `diana@example.com` | member | — | member | lead | — | member |
| 5 | Eve Developer | `eve@example.com` | member | — | member | member | — | member |
| 6 | Frank Designer | `frank@example.com` | member | — | — | member | — | — |
| 7 | Grace Tester | `grace@example.com` | member | — | — | member | — | — |
| 8 | Hank Analyst | `hank@example.com` | — | admin | — | — | admin | — |
| 9 | Ivy Intern | `ivy@example.com` | — | member | — | — | member | — |
| 10 | Jack Viewer | `jack@example.com` | — | — | admin | — | — | admin |

> **Primary test account**: `test@example.com` / `password123` — has owner access to everything.

---

## API Reference

All endpoints except authentication and health check require the `Authorization: Bearer <token>` header.

Base URL: `http://localhost:8080`

### Authentication (Public)

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/auth/register` | Register a new user |
| `POST` | `/auth/login` | Login, returns JWT token + user |

### Health

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Returns `{"status":"ok"}` |

### Invitations (User-Level)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/invitations` | List my pending workspace invitations |
| `PATCH` | `/invitations/{invitationID}` | Accept or decline a workspace invitation |
| `GET` | `/org-invitations` | List my pending org invitations |
| `PATCH` | `/org-invitations/{invitationID}` | Accept or decline an org invitation |

### Organizations

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/organizations` | List my organizations |
| `POST` | `/organizations` | Create a new organization |
| `GET` | `/organizations/{orgID}` | Get organization details |
| `PATCH` | `/organizations/{orgID}` | Update organization |
| `GET` | `/organizations/{orgID}/members` | List org members |
| `DELETE` | `/organizations/{orgID}/members/{userID}` | Remove a member from org (cascades) |
| `POST` | `/organizations/{orgID}/members/leave` | Leave the organization (cascades) |
| `GET` | `/organizations/{orgID}/prefixes` | List all project prefixes in org |
| `GET` | `/organizations/{orgID}/tasks/by-key/{taskKey}` | Get task by Jira-style key (e.g. `APIGW-001`) |

### Org Invitations

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/organizations/{orgID}/invitations` | Invite a user to the org |
| `GET` | `/organizations/{orgID}/invitations` | List org invitations |

### Teams

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/organizations/{orgID}/teams` | List teams in org |
| `POST` | `/organizations/{orgID}/teams` | Create a team (`manager`+ required) |
| `GET` | `/organizations/{orgID}/teams/{teamID}` | Get team details |
| `PATCH` | `/organizations/{orgID}/teams/{teamID}` | Update team |
| `DELETE` | `/organizations/{orgID}/teams/{teamID}` | Delete team |
| `POST` | `/organizations/{orgID}/teams/{teamID}/members` | Add member to team |
| `DELETE` | `/organizations/{orgID}/teams/{teamID}/members/{userID}` | Remove member from team |

### Workspaces

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/organizations/{orgID}/workspaces` | List workspaces in org |
| `POST` | `/organizations/{orgID}/workspaces` | Create workspace |
| `GET` | `/organizations/{orgID}/workspaces/{workspaceID}` | Get workspace details |
| `PATCH` | `/organizations/{orgID}/workspaces/{workspaceID}` | Update workspace |
| `DELETE` | `/organizations/{orgID}/workspaces/{workspaceID}` | Delete workspace (owner only) |
| `GET` | `/organizations/{orgID}/workspaces/{workspaceID}/stats` | Dashboard statistics |

### Workspace Members

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `.../workspaces/{workspaceID}/members` | List workspace members |
| `POST` | `.../workspaces/{workspaceID}/members` | Invite member by email |
| `POST` | `.../workspaces/{workspaceID}/members/add` | Direct-add an org member (with explicit role) |
| `POST` | `.../workspaces/{workspaceID}/members/leave` | Leave workspace (cascades to projects) |
| `PATCH` | `.../workspaces/{workspaceID}/members/{userID}` | Change member role |
| `DELETE` | `.../workspaces/{workspaceID}/members/{userID}` | Remove member |
| `GET` | `.../workspaces/{workspaceID}/invitations` | List pending invitations |

### Workspace Teams

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `.../workspaces/{workspaceID}/teams` | Assign a team to workspace |
| `DELETE` | `.../workspaces/{workspaceID}/teams/{teamID}` | Remove a team from workspace |

### Projects

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `.../workspaces/{workspaceID}/projects` | List projects (`?page=`, `?limit=`, `?search=`, `?owner=`) |
| `POST` | `.../workspaces/{workspaceID}/projects` | Create project (with unique prefix) |
| `GET` | `/projects/{projectID}` | Get project details |
| `PATCH` | `/projects/{projectID}` | Update project |
| `DELETE` | `/projects/{projectID}` | Delete project |

### Project Members

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `.../projects/{projectID}/members` | List project members |
| `POST` | `.../projects/{projectID}/members` | Add member (must be workspace member, explicit role) |
| `POST` | `.../projects/{projectID}/members/leave` | Leave project |
| `PATCH` | `.../projects/{projectID}/members/{userID}` | Change member role |
| `DELETE` | `.../projects/{projectID}/members/{userID}` | Remove member |
| `POST` | `.../projects/{projectID}/teams` | Assign a team to project |

### Tasks

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/projects/{projectID}/tasks` | List tasks (`?status=`, `?assignee=`, `?priority=`, `?search=`, `?page=`, `?limit=`) |
| `POST` | `/projects/{projectID}/tasks` | Create task (auto-generates task key) |
| `PATCH` | `/tasks/{taskID}` | Update task |
| `DELETE` | `/tasks/{taskID}` | Delete task |

### Comments

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/projects/{projectID}/comments` | List project-level comments |
| `POST` | `/projects/{projectID}/comments` | Create project comment |
| `GET` | `/tasks/{taskID}/comments` | List task comments |
| `POST` | `/tasks/{taskID}/comments` | Create task comment |
| `PATCH` | `/comments/{commentID}` | Edit comment (author only) |
| `DELETE` | `/comments/{commentID}` | Delete comment (author only) |

### Custom Fields

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/projects/{projectID}/custom-fields` | List custom field definitions |
| `POST` | `/projects/{projectID}/custom-fields` | Create custom field definition |
| `DELETE` | `/projects/{projectID}/custom-fields/{fieldID}` | Delete custom field definition |
| `GET` | `/tasks/{taskID}/custom-fields` | Get field values for a task |
| `PUT` | `/tasks/{taskID}/custom-fields/{fieldID}` | Set field value on a task |

### WIP Limits

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/projects/{projectID}/wip-limits` | Get WIP limits for project |
| `PUT` | `/projects/{projectID}/wip-limits` | Set WIP limit for a status column |
| `DELETE` | `/projects/{projectID}/wip-limits` | Remove WIP limit |

### Personal Dashboard

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/dashboard/my-tasks` | Paginated tasks assigned to the current user (`?status=`, `?priority=`, `?search=`, `?project_id=`, `?due_before=`, `?due_after=`, `?page=`, `?limit=`) |
| `GET` | `/dashboard/my-stats` | Aggregated stats: total, by status, by priority, overdue, due today, due this week, completed |
| `GET` | `/dashboard/my-projects` | List of projects the user has tasks in (id + name) |

### Org Admin Dashboard (manager+ only)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/organizations/{orgID}/dashboard/member-stats` | Task stats for all org members below the caller's role level |
| `GET` | `/organizations/{orgID}/dashboard/member-tasks/{userID}` | Paginated tasks for a specific member (same filters as my-tasks) |

### WebSocket (Real-Time)

| Endpoint | Auth | Description |
|---|---|---|
| `GET /ws?token=<jwt>` | Via query param | Upgrades to WebSocket. Receives real-time events. |

**Events pushed to clients:**

| Event Type | Payload | When |
|---|---|---|
| `task_created` | Full task object + `project_id` | A task is created anywhere |
| `task_updated` | Full task object + `project_id` | A task is updated (status, assignee, etc.) |
| `task_deleted` | `{"id": "<task-uuid>"}` | A task is deleted |
| `force_logout` | `{"reason": "token_expired"}` | JWT token has expired — client should log out |

The WebSocket connection uses ping/pong keepalive (54s interval) and checks token expiry every 30 seconds. If the token is expired, the server sends `force_logout` and closes the connection.

### Error Responses

All errors follow a consistent JSON format:

```json
{ "error": "validation failed", "fields": { "email": "is required" } }
```

| Status Code | Meaning |
|---|---|
| `400` | Validation error or bad request |
| `401` | Missing or invalid JWT token |
| `403` | Insufficient permissions |
| `404` | Resource not found |
| `409` | Conflict (duplicate email, slug, etc.) |

---

## Project Structure

```
TaskFlow/
├── docker-compose.yml
├── .env.example
│
├── backend/
│   ├── Dockerfile
│   ├── entrypoint.sh              # Migrations → Seed → Server
│   ├── go.mod
│   ├── cmd/server/main.go         # App entry point, DI wiring, routes
│   ├── scripts/seed/main.go       # Go seeder (runs on every startup)
│   ├── migrations/                # 14 sequential SQL migrations
│   │   ├── 000001_create_users
│   │   ├── 000002_create_workspaces
│   │   ├── 000003_create_workspace_members
│   │   ├── 000004_create_projects
│   │   ├── 000005_create_tasks
│   │   ├── 000006_expand_roles
│   │   ├── 000007_create_invitations
│   │   ├── 000008_add_start_date
│   │   ├── 000009_create_wip_limits
│   │   ├── 000010_create_custom_fields
│   │   ├── 000011_create_comments
│   │   ├── 000012_add_project_prefix_and_task_numbering
│   │   ├── 000013_create_organizations
│   │   └── 000014_teams_and_project_members
│   └── internal/
│       ├── config/                # Environment-based configuration
│       ├── domain/                # Core entities (User, Org, Workspace, Project, Task, Team, etc.)
│       ├── dto/                   # Request/response data transfer objects
│       ├── handler/               # HTTP handlers (one per domain)
│       ├── middleware/            # Auth, OrgGuard, WorkspaceGuard, Logger, Recovery
│       ├── repository/           # PostgreSQL queries via pgx
│       ├── service/              # Business logic layer
│       └── ws/                   # WebSocket hub + client (real-time events)
│
└── frontend/
    ├── Dockerfile
    ├── nginx.conf
    ├── package.json
    ├── vite.config.ts
    └── src/
        ├── App.tsx                # Route definitions
        ├── components/
        │   ├── layout/            # Navbar, Sidebar
        │   └── ui/                # Button, Pagination, etc.
        ├── core/
        │   ├── api/               # Axios HTTP clients per domain (with global error toasts)
        │   ├── guards/            # AuthGuard (token expiry + WebSocket init)
        │   └── providers/         # QueryProvider, ThemeProvider
        ├── features/
        │   ├── tasks/             # KanbanBoard, TaskCard, TaskModal
        │   ├── workspaces/        # MemberList, WorkspaceCard
        │   └── projects/          # ProjectMemberList
        ├── hooks/                 # TanStack Query hooks per domain + useWebSocket
        ├── pages/                 # Route-level page components (incl. MyDashboardPage)
        ├── store/                 # Zustand auth store
        ├── types/                 # TypeScript interfaces
        └── utils/                 # Helpers (cn, formatDate, etc.)
```

---

## Environment Variables

Copy `.env.example` to `.env` before first run. For production, point `DATABASE_URL` to a [Neon](https://neon.tech) serverless Postgres instance (or any PostgreSQL 16+ provider).

| Variable | Default | Description |
|---|---|---|
| `POSTGRES_USER` | `taskflow` | PostgreSQL username (local Docker) |
| `POSTGRES_PASSWORD` | `taskflow_secret` | PostgreSQL password (local Docker) |
| `POSTGRES_DB` | `taskflow` | PostgreSQL database name (local Docker) |
| `DATABASE_URL` | `postgres://taskflow:taskflow_secret@db:5432/taskflow?sslmode=disable` | Full connection string — use a Neon connection string for production |
| `JWT_SECRET` | `change-me-to-a-random-secret` | Secret for signing JWT tokens |
| `API_PORT` | `8080` | Backend API port |
| `FRONTEND_PORT` | `3000` | Frontend port |
| `BCRYPT_COST` | `12` | bcrypt hashing cost factor |

---

## What You'd Do With More Time

### Next Steps

- **Database-backed integration tests** — Spin up a test Postgres container and run full HTTP round-trip tests (register → create org → create workspace → create task → verify). Current tests cover validation and JWT logic but not the full request lifecycle.
- **Task ordering** — Persist sort order within Kanban columns (currently sorted by creation date).
- **File attachments** — Attach files to tasks with S3-compatible storage.
- **Activity feed** — Audit log of all changes on tasks, members, and settings.
- **Email invitations** — Send actual invitation emails via SMTP (currently in-app only).
- **Rate limiting** — Per-IP and per-user rate limiting on auth endpoints.
- **E2E tests** — Playwright tests for critical user flows (login, create task, drag-and-drop, WebSocket reconnection).
- **CI/CD pipeline** — GitHub Actions for linting, testing, and building Docker images.
- **Notifications** — In-app notification center for mentions, assignments, and invitation responses.
- **Search** — Global search across tasks, projects, and members with keyboard shortcut.
