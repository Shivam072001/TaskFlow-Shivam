import type React from 'react';
import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Loader2, Mail, Lock, Eye, EyeOff, UserRound, LayoutDashboard, Check, X } from 'lucide-react';
import { cn } from '@utils/cn';

interface Props {
  onSubmit: (name: string, email: string, password: string) => void;
  isLoading: boolean;
  error: string | null;
}

export function RegisterForm({ onSubmit, isLoading, error }: Props) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  const passwordChecks = useMemo(() => [
    { label: 'At least 6 characters', met: password.length >= 6 },
    { label: 'Contains a letter', met: /[a-zA-Z]/.test(password) },
    { label: 'Contains a number', met: /\d/.test(password) },
  ], [password]);

  const passwordStrength = passwordChecks.filter((c) => c.met).length;

  const validate = () => {
    const e: Record<string, string> = {};
    if (!name.trim()) e.name = 'Name is required';
    if (!email.trim()) e.email = 'Email is required';
    else if (!/\S+@\S+\.\S+/.test(email)) e.email = 'Invalid email format';
    if (password.length < 6) e.password = 'Must be at least 6 characters';
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleBlur = (field: string) => {
    setTouched((p) => ({ ...p, [field]: true }));
    validate();
  };

  const handleSubmit = (ev: React.FormEvent<HTMLFormElement>) => {
    ev.preventDefault();
    setTouched({ name: true, email: true, password: true });
    if (validate()) onSubmit(name, email, password);
  };

  const isValid = name.trim() && email.trim() && /\S+@\S+\.\S+/.test(email) && password.length >= 6;

  return (
    <div className="w-full max-w-md">
      <div className="mb-8 text-center">
        <Link to="/" className="inline-flex items-center gap-2.5 mb-8">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary text-primary-foreground">
            <LayoutDashboard className="h-5 w-5" />
          </div>
          <span className="text-xl font-bold text-foreground">TaskFlow</span>
        </Link>
        <h1 className="text-2xl font-bold text-foreground">Create your account</h1>
        <p className="mt-2 text-muted-foreground">
          Start managing your projects in minutes
        </p>
      </div>

      <div className="rounded-2xl border border-border bg-card p-8 shadow-sm">
        {error && (
          <div className="mb-6 flex items-center gap-2 rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
            <div className="h-1.5 w-1.5 rounded-full bg-destructive" />
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="mb-2 block text-sm font-medium text-foreground">
              Full name
            </label>
            <div className="relative">
              <UserRound className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                onBlur={() => handleBlur('name')}
                className={cn(
                  'w-full rounded-xl border bg-background py-3 pl-10 pr-4 text-sm text-foreground',
                  'placeholder:text-muted-foreground transition-all',
                  'focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary',
                  touched.name && errors.name
                    ? 'border-destructive focus:ring-destructive/50'
                    : 'border-border',
                )}
                placeholder="John Doe"
                autoComplete="name"
                autoFocus
              />
            </div>
            {touched.name && errors.name && (
              <p className="mt-1.5 text-xs text-destructive">{errors.name}</p>
            )}
          </div>

          <div>
            <label className="mb-2 block text-sm font-medium text-foreground">
              Email address
            </label>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                onBlur={() => handleBlur('email')}
                className={cn(
                  'w-full rounded-xl border bg-background py-3 pl-10 pr-4 text-sm text-foreground',
                  'placeholder:text-muted-foreground transition-all',
                  'focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary',
                  touched.email && errors.email
                    ? 'border-destructive focus:ring-destructive/50'
                    : 'border-border',
                )}
                placeholder="you@example.com"
                autoComplete="email"
              />
            </div>
            {touched.email && errors.email && (
              <p className="mt-1.5 text-xs text-destructive">{errors.email}</p>
            )}
          </div>

          <div>
            <label className="mb-2 block text-sm font-medium text-foreground">
              Password
            </label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                onBlur={() => handleBlur('password')}
                className={cn(
                  'w-full rounded-xl border bg-background py-3 pl-10 pr-12 text-sm text-foreground',
                  'placeholder:text-muted-foreground transition-all',
                  'focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary',
                  touched.password && errors.password
                    ? 'border-destructive focus:ring-destructive/50'
                    : 'border-border',
                )}
                placeholder="Create a password"
                autoComplete="new-password"
              />
              <button
                type="button"
                onClick={() => setShowPassword((p) => !p)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                tabIndex={-1}
              >
                {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>

            {password && (
              <div className="mt-3 space-y-2">
                <div className="flex gap-1">
                  {[1, 2, 3].map((i) => (
                    <div
                      key={i}
                      className={cn(
                        'h-1 flex-1 rounded-full transition-colors',
                        i <= passwordStrength
                          ? passwordStrength === 1
                            ? 'bg-red-500'
                            : passwordStrength === 2
                              ? 'bg-amber-500'
                              : 'bg-green-500'
                          : 'bg-border',
                      )}
                    />
                  ))}
                </div>
                <div className="space-y-1">
                  {passwordChecks.map((check) => (
                    <div
                      key={check.label}
                      className={cn(
                        'flex items-center gap-1.5 text-xs transition-colors',
                        check.met ? 'text-green-600 dark:text-green-400' : 'text-muted-foreground',
                      )}
                    >
                      {check.met ? (
                        <Check className="h-3 w-3" />
                      ) : (
                        <X className="h-3 w-3" />
                      )}
                      {check.label}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          <button
            type="submit"
            disabled={isLoading || !isValid}
            className={cn(
              'relative w-full rounded-xl bg-primary py-3 text-sm font-semibold text-primary-foreground',
              'hover:bg-primary/90 transition-all disabled:opacity-50',
              'shadow-sm hover:shadow-md hover:shadow-primary/20',
            )}
          >
            {isLoading ? (
              <Loader2 className="mx-auto h-5 w-5 animate-spin" />
            ) : (
              'Create account'
            )}
          </button>
        </form>
      </div>

      <p className="mt-6 text-center text-sm text-muted-foreground">
        Already have an account?{' '}
        <Link to="/login" className="font-medium text-primary hover:text-primary/80 transition-colors">
          Sign in
        </Link>
      </p>
    </div>
  );
}
