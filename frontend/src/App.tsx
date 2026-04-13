import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryProvider } from '@core/providers/QueryProvider';
import { ThemeProvider } from '@core/providers/ThemeProvider';
import { AuthGuard } from '@core/guards/AuthGuard';
import { Toaster } from '@components/ui/Toaster';
import { ConfirmDialogProvider } from '@components/ui/ConfirmDialog';
import { useAuth } from '@hooks/useAuth';
import { LandingPage } from '@pages/LandingPage';
import { LoginPage } from '@pages/LoginPage';
import { RegisterPage } from '@pages/RegisterPage';
import { OrgSelectorPage } from '@pages/OrgSelectorPage';
import { WorkspaceListPage } from '@pages/WorkspaceListPage';
import { WorkspaceDashboardPage } from '@pages/WorkspaceDashboardPage';
import { ProjectDetailPage } from '@pages/ProjectDetailPage';
import { TaskKeyPage } from '@pages/TaskKeyPage';
import { OrgSettingsPage } from '@pages/OrgSettingsPage';
import { MyDashboardPage } from '@pages/MyDashboardPage';
import type { ReactNode } from 'react';

function ProtectedRoute({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

function PublicRoute({ children }: { children: ReactNode }) {
  const { isAuthenticated } = useAuth();
  if (isAuthenticated) return <Navigate to="/organizations" replace />;
  return <>{children}</>;
}

export default function App() {
  return (
    <QueryProvider>
      <ThemeProvider>
        <ConfirmDialogProvider>
        <BrowserRouter>
          <AuthGuard>
            <Routes>
              <Route path="/" element={<PublicRoute><LandingPage /></PublicRoute>} />
              <Route path="/login" element={<PublicRoute><LoginPage /></PublicRoute>} />
              <Route path="/register" element={<PublicRoute><RegisterPage /></PublicRoute>} />
              <Route path="/dashboard" element={<ProtectedRoute><MyDashboardPage /></ProtectedRoute>} />
              <Route path="/organizations" element={<ProtectedRoute><OrgSelectorPage /></ProtectedRoute>} />
              <Route path="/org/:orgSlug/workspaces" element={<ProtectedRoute><WorkspaceListPage /></ProtectedRoute>} />
              <Route path="/org/:orgSlug/workspaces/:wid" element={<ProtectedRoute><WorkspaceDashboardPage /></ProtectedRoute>} />
              <Route path="/org/:orgSlug/workspaces/:wid/projects/:pid" element={<ProtectedRoute><ProjectDetailPage /></ProtectedRoute>} />
              <Route path="/org/:orgSlug/task/:taskKey" element={<ProtectedRoute><TaskKeyPage /></ProtectedRoute>} />
              <Route path="/org/:orgSlug/settings" element={<ProtectedRoute><OrgSettingsPage /></ProtectedRoute>} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </AuthGuard>
        </BrowserRouter>
        <Toaster />
        </ConfirmDialogProvider>
      </ThemeProvider>
    </QueryProvider>
  );
}
