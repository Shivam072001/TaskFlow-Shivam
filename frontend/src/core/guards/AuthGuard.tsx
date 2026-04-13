import { useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@hooks/useAuth';
import { useWebSocket } from '@hooks/useWebSocket';
import { isTokenExpired } from '@utils/jwt';

const CHECK_INTERVAL_MS = 60_000;

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const { token, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const hasLoggedOut = useRef(false);

  useWebSocket();

  useEffect(() => {
    hasLoggedOut.current = false;
  }, [token]);

  useEffect(() => {
    if (!isAuthenticated || !token) return;

    const handleExpiry = () => {
      if (hasLoggedOut.current) return;
      if (isTokenExpired(token)) {
        hasLoggedOut.current = true;
        logout();
        navigate('/login', { replace: true });
      }
    };

    handleExpiry();

    const id = setInterval(handleExpiry, CHECK_INTERVAL_MS);
    return () => clearInterval(id);
  }, [isAuthenticated, token, logout, navigate]);

  return <>{children}</>;
}
