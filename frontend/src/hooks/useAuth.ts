import { create } from 'zustand';
import { isTokenExpired } from '@utils/jwt';
import { queryClient } from '@core/providers/queryClient';
import type { Organization } from '@/types';

interface AuthUser {
  id: string;
  name: string;
  email: string;
}

interface AuthState {
  token: string | null;
  user: AuthUser | null;
  isAuthenticated: boolean;
  activeOrg: Organization | null;
  login: (token: string, user: AuthUser) => void;
  logout: () => void;
  setActiveOrg: (org: Organization) => void;
  clearActiveOrg: () => void;
}

function safeParse<T>(raw: string | null): T | null {
  if (!raw) return null;
  try {
    return JSON.parse(raw) as T;
  } catch {
    return null;
  }
}

export const useAuth = create<AuthState>((set) => {
  const token = localStorage.getItem('token');
  const user = safeParse<AuthUser>(localStorage.getItem('user'));
  const activeOrg = safeParse<Organization>(localStorage.getItem('activeOrg'));

  const expired = token ? isTokenExpired(token) : false;
  if (expired) {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('activeOrg');
  }

  return {
    token: expired ? null : token,
    user: expired ? null : user,
    isAuthenticated: !expired && !!token && !!user,
    activeOrg: expired ? null : activeOrg,
    login: (token: string, user: AuthUser) => {
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      set({ token, user, isAuthenticated: true });
    },
    logout: () => {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      localStorage.removeItem('activeOrg');
      queryClient.clear();
      set({ token: null, user: null, isAuthenticated: false, activeOrg: null });
    },
    setActiveOrg: (org: Organization) => {
      localStorage.setItem('activeOrg', JSON.stringify(org));
      set({ activeOrg: org });
    },
    clearActiveOrg: () => {
      localStorage.removeItem('activeOrg');
      set({ activeOrg: null });
    },
  };
});
