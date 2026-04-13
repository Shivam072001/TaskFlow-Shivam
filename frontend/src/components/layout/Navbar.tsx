import { Link, useNavigate, useLocation } from 'react-router-dom';
import {
  LogOut, Moon, Sun, LayoutDashboard, Building2,
  Settings, User, LogIn, UserPlus,
} from 'lucide-react';
import { useAuth } from '@hooks/useAuth';
import { useTheme } from '@hooks/useTheme';
import { cn } from '@utils/cn';

export function Navbar() {
  const { user, isAuthenticated, activeOrg, logout, clearActiveOrg } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();
  const { pathname } = useLocation();

  const homeLink = isAuthenticated
    ? activeOrg ? `/org/${activeOrg.slug}/workspaces` : '/organizations'
    : '/';

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleSwitchOrg = () => {
    clearActiveOrg();
    navigate('/organizations');
  };

  const isLanding = pathname === '/';

  return (
    <header
      className={cn(
        'sticky top-0 z-50 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 transition-colors',
        isLanding ? 'border-transparent' : 'border-border',
      )}
    >
      <div className="mx-auto flex h-16 max-w-7xl items-center px-4 md:px-6">
        <Link to={homeLink} className="flex items-center gap-2.5 font-bold text-foreground text-lg">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <LayoutDashboard className="h-4.5 w-4.5" />
          </div>
          <span>TaskFlow</span>
        </Link>

        {isAuthenticated && activeOrg && (
          <>
            <span className="mx-2 text-border">/</span>
            <button
              onClick={handleSwitchOrg}
              className="flex items-center gap-1.5 rounded-md border border-border px-2.5 py-1 text-xs text-muted-foreground hover:bg-muted transition-colors"
            >
              <Building2 className="h-3.5 w-3.5" />
              <span className="max-w-[120px] truncate">{activeOrg.name}</span>
            </button>
            <Link
              to={`/org/${activeOrg.slug}/settings`}
              className="ml-1 inline-flex h-7 w-7 items-center justify-center rounded-md text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
              title="Organization Settings"
            >
              <Settings className="h-3.5 w-3.5" />
            </Link>
          </>
        )}

        <div className="ml-auto flex items-center gap-2">
          {isAuthenticated ? (
            <>
              <Link
                to="/dashboard"
                className={cn(
                  'inline-flex h-9 items-center gap-1.5 rounded-md px-3',
                  'text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors',
                )}
                title="My Dashboard"
              >
                <User className="h-4 w-4" />
                <span className="hidden md:inline">My Tasks</span>
              </Link>

              <button
                onClick={toggleTheme}
                aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
                className="inline-flex h-9 w-9 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
              </button>

              <div className="ml-1 flex items-center gap-2 border-l border-border pl-3">
                <div className="hidden items-center gap-2 md:flex">
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-sm font-medium text-primary">
                    {user?.name?.charAt(0).toUpperCase()}
                  </div>
                  <span className="text-sm text-muted-foreground">{user?.name}</span>
                </div>
                <button
                  onClick={handleLogout}
                  aria-label="Log out"
                  className="inline-flex h-9 w-9 items-center justify-center rounded-md text-muted-foreground hover:bg-destructive/10 hover:text-destructive transition-colors"
                >
                  <LogOut className="h-4 w-4" />
                </button>
              </div>
            </>
          ) : (
            <>
              <button
                onClick={toggleTheme}
                aria-label={theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
                className="inline-flex h-9 w-9 items-center justify-center rounded-md text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                {theme === 'dark' ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
              </button>

              <Link
                to="/login"
                className="inline-flex h-9 items-center gap-1.5 rounded-md px-4 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                <LogIn className="h-4 w-4" />
                Sign in
              </Link>
              <Link
                to="/register"
                className="inline-flex h-9 items-center gap-1.5 rounded-lg bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
              >
                <UserPlus className="h-4 w-4" />
                Sign up
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
