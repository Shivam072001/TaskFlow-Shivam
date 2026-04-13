import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { RegisterForm } from '@features/auth/RegisterForm';
import { Navbar } from '@components/layout/Navbar';
import { useAuth } from '@hooks/useAuth';
import * as authApi from '@core/api/auth';
import type { ApiError } from '@/types';
import { AxiosError } from 'axios';
import { Shield, Layers, Zap } from 'lucide-react';

const perks = [
  { icon: Zap, title: 'Instant setup', desc: 'Create your organization and start tracking in under a minute.' },
  { icon: Layers, title: 'Flexible hierarchy', desc: 'Orgs, workspaces, projects — structure that scales with you.' },
  { icon: Shield, title: 'Secure by default', desc: 'JWT auth, role-based access, and encrypted passwords.' },
];

export function RegisterPage() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (name: string, email: string, password: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await authApi.register(name, email, password);
      login(res.token, res.user);
      navigate('/organizations');
    } catch (err) {
      const axiosErr = err as AxiosError<ApiError>;
      const fields = axiosErr.response?.data?.fields;
      setError(fields ? Object.values(fields).join(', ') : (axiosErr.response?.data?.error ?? 'Registration failed'));
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
              Your team's new home for work
            </h2>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              Join thousands of teams who use TaskFlow to organize, prioritize, and deliver work on time.
            </p>
            <div className="mt-10 space-y-6">
              {perks.map((p) => (
                <div key={p.title} className="flex gap-4">
                  <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary">
                    <p.icon className="h-5 w-5" />
                  </div>
                  <div>
                    <div className="text-sm font-semibold text-foreground">{p.title}</div>
                    <div className="mt-0.5 text-sm text-muted-foreground">{p.desc}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Right panel — form */}
        <div className="flex flex-1 items-center justify-center p-6">
          <RegisterForm onSubmit={handleSubmit} isLoading={isLoading} error={error} />
        </div>
      </div>
    </div>
  );
}
