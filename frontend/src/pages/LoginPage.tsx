import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { LoginForm } from '@features/auth/LoginForm';
import { Navbar } from '@components/layout/Navbar';
import { useAuth } from '@hooks/useAuth';
import * as authApi from '@core/api/auth';
import type { ApiError } from '@/types';
import { AxiosError } from 'axios';
import { CheckCircle2 } from 'lucide-react';

const highlights = [
  'Multi-tenant organizations & workspaces',
  'Drag-and-drop Kanban boards',
  'Real-time WebSocket updates',
  'Granular role-based access control',
  'Personal dashboards & analytics',
];

export function LoginPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (email: string, password: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await authApi.login(email, password);
      login(res.token, res.user);
      navigate('/organizations');
    } catch (err) {
      const axiosErr = err as AxiosError<ApiError>;
      setError(axiosErr.response?.data?.error ?? 'Login failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <Navbar />
      <div className="flex flex-1">
        {/* Left panel — branding (hidden on mobile) */}
        <div className="hidden w-1/2 flex-col justify-center bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5 p-12 lg:flex xl:p-20">
          <div className="max-w-md">
            <h2 className="text-3xl font-bold tracking-tight text-foreground">
              Manage projects with clarity
            </h2>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              TaskFlow gives your team the structure, visibility, and speed to ship great work — every single day.
            </p>
            <ul className="mt-8 space-y-3">
              {highlights.map((h) => (
                <li key={h} className="flex items-center gap-3 text-sm text-foreground">
                  <CheckCircle2 className="h-4 w-4 shrink-0 text-primary" />
                  {h}
                </li>
              ))}
            </ul>
          </div>
        </div>

        {/* Right panel — form */}
        <div className="flex flex-1 items-center justify-center p-6">
          <LoginForm onSubmit={handleSubmit} isLoading={isLoading} error={error} />
        </div>
      </div>
    </div>
  );
}
