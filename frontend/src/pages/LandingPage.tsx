import { Link } from 'react-router-dom';
import {
  ArrowRight, CheckCircle2, Zap, Shield, Users, Layers,
  BarChart3, Globe, Kanban,
} from 'lucide-react';
import { Navbar } from '@components/layout/Navbar';
import { cn } from '@utils/cn';

const features = [
  {
    icon: Kanban,
    title: 'Kanban Boards',
    description: 'Drag-and-drop task boards with real-time updates. Organize work visually across customizable status columns.',
  },
  {
    icon: Users,
    title: 'Team Collaboration',
    description: 'Invite members, assign roles at every level, and work together seamlessly across organizations and workspaces.',
  },
  {
    icon: Shield,
    title: 'Granular Permissions',
    description: 'Three-tier role system across organizations, workspaces, and projects. Control who sees and does what.',
  },
  {
    icon: Layers,
    title: 'Multi-Tenancy',
    description: 'Organizations, workspaces, and projects — a flexible hierarchy that scales from solo use to enterprise teams.',
  },
  {
    icon: BarChart3,
    title: 'Dashboards & Stats',
    description: 'Personal task dashboards, project statistics, and org-level member analytics — all at a glance.',
  },
  {
    icon: Globe,
    title: 'Real-Time Updates',
    description: 'WebSocket-powered live updates. See task changes the instant they happen, no refresh needed.',
  },
];

const stats = [
  { value: '18', label: 'Database tables' },
  { value: '73', label: 'API endpoints' },
  { value: '7', label: 'Role levels' },
  { value: '∞', label: 'Possibilities' },
];

export function LandingPage() {
  return (
    <div className="min-h-screen bg-background text-foreground">
      <Navbar />

      {/* Hero */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 -z-10">
          <div className="absolute left-1/2 top-0 -translate-x-1/2 h-[600px] w-[800px] rounded-full bg-primary/5 blur-3xl" />
          <div className="absolute right-0 top-1/4 h-[400px] w-[400px] rounded-full bg-primary/3 blur-3xl" />
        </div>

        <div className="mx-auto max-w-7xl px-4 pb-20 pt-20 md:px-6 md:pt-32 md:pb-28">
          <div className="mx-auto max-w-3xl text-center">
            <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-primary/20 bg-primary/5 px-4 py-1.5 text-sm font-medium text-primary">
              <Zap className="h-3.5 w-3.5" />
              Built with Go + React + PostgreSQL
            </div>

            <h1 className="text-4xl font-extrabold tracking-tight sm:text-5xl md:text-6xl">
              Task management{' '}
              <span className="bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                that scales
              </span>{' '}
              with your team
            </h1>

            <p className="mx-auto mt-6 max-w-2xl text-lg text-muted-foreground leading-relaxed">
              TaskFlow is a full-stack, multi-tenant project management platform with
              organizations, workspaces, drag-and-drop boards, real-time collaboration,
              and granular role-based access control.
            </p>

            <div className="mt-10 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
              <Link
                to="/register"
                className={cn(
                  'inline-flex h-12 items-center gap-2 rounded-xl bg-primary px-8 text-base font-semibold text-primary-foreground',
                  'shadow-lg shadow-primary/25 hover:bg-primary/90 hover:shadow-xl hover:shadow-primary/30 transition-all',
                )}
              >
                Get Started Free
                <ArrowRight className="h-4 w-4" />
              </Link>
              <Link
                to="/login"
                className="inline-flex h-12 items-center gap-2 rounded-xl border border-border bg-card px-8 text-base font-semibold text-foreground hover:bg-accent transition-colors"
              >
                Sign in
              </Link>
            </div>
          </div>

          {/* Stats strip */}
          <div className="mx-auto mt-20 grid max-w-2xl grid-cols-2 gap-6 sm:grid-cols-4">
            {stats.map((s) => (
              <div key={s.label} className="text-center">
                <div className="text-3xl font-extrabold text-foreground">{s.value}</div>
                <div className="mt-1 text-sm text-muted-foreground">{s.label}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="border-t border-border bg-muted/30 py-20 md:py-28">
        <div className="mx-auto max-w-7xl px-4 md:px-6">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              Everything you need to ship faster
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              A complete platform built from the ground up — no compromises.
            </p>
          </div>

          <div className="mx-auto mt-16 grid max-w-5xl gap-8 sm:grid-cols-2 lg:grid-cols-3">
            {features.map((f) => (
              <div
                key={f.title}
                className="group rounded-2xl border border-border bg-card p-6 transition-all hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5"
              >
                <div className="mb-4 inline-flex h-11 w-11 items-center justify-center rounded-xl bg-primary/10 text-primary transition-colors group-hover:bg-primary group-hover:text-primary-foreground">
                  <f.icon className="h-5 w-5" />
                </div>
                <h3 className="text-lg font-semibold text-foreground">{f.title}</h3>
                <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                  {f.description}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How it works */}
      <section className="py-20 md:py-28">
        <div className="mx-auto max-w-7xl px-4 md:px-6">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              Up and running in minutes
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              Three steps to organized teamwork.
            </p>
          </div>

          <div className="mx-auto mt-16 grid max-w-4xl gap-8 md:grid-cols-3">
            {[
              { step: '01', title: 'Create your organization', desc: 'Set up your org and invite your team. Everyone gets the right level of access from day one.' },
              { step: '02', title: 'Organize into workspaces', desc: 'Group related projects into workspaces. Each workspace has its own members and settings.' },
              { step: '03', title: 'Track and deliver', desc: 'Create tasks, drag them across your board, and watch progress update in real time.' },
            ].map((item) => (
              <div key={item.step} className="relative text-center">
                <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-primary/10 text-xl font-extrabold text-primary">
                  {item.step}
                </div>
                <h3 className="text-lg font-semibold text-foreground">{item.title}</h3>
                <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{item.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Tech Stack */}
      <section className="border-t border-border bg-muted/30 py-20 md:py-28">
        <div className="mx-auto max-w-7xl px-4 md:px-6">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              Built with modern tools
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              Production-grade stack from frontend to infrastructure.
            </p>
          </div>

          <div className="mx-auto mt-12 grid max-w-3xl grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4">
            {[
              'Go (Chi)', 'React 19', 'TypeScript', 'PostgreSQL 16',
              'Tailwind CSS 4', 'TanStack Query', 'WebSocket', 'Docker',
            ].map((tech) => (
              <div
                key={tech}
                className="flex items-center justify-center gap-2 rounded-xl border border-border bg-card px-4 py-3 text-sm font-medium text-foreground transition-colors hover:border-primary/30"
              >
                <CheckCircle2 className="h-3.5 w-3.5 text-primary" />
                {tech}
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-20 md:py-28">
        <div className="mx-auto max-w-7xl px-4 md:px-6">
          <div className="mx-auto max-w-2xl rounded-3xl border border-primary/20 bg-gradient-to-br from-primary/5 to-primary/10 p-10 text-center md:p-16">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              Ready to get organized?
            </h2>
            <p className="mx-auto mt-4 max-w-md text-muted-foreground">
              Start managing your projects with TaskFlow today. Free to use, easy to deploy.
            </p>
            <div className="mt-8 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
              <Link
                to="/register"
                className={cn(
                  'inline-flex h-12 items-center gap-2 rounded-xl bg-primary px-8 text-base font-semibold text-primary-foreground',
                  'shadow-lg shadow-primary/25 hover:bg-primary/90 transition-all',
                )}
              >
                Create Free Account
                <ArrowRight className="h-4 w-4" />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border py-8">
        <div className="mx-auto flex max-w-7xl flex-col items-center gap-4 px-4 md:flex-row md:justify-between md:px-6">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <Kanban className="h-3.5 w-3.5" />
            </div>
            TaskFlow
          </div>
          <p className="text-xs text-muted-foreground">
            Built as a full-stack engineering showcase. Open source.
          </p>
        </div>
      </footer>
    </div>
  );
}
