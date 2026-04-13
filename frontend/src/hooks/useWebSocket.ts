import { useEffect, useRef, useCallback } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { useAuth } from './useAuth';
import { WS_BASE_URL } from '@core/config';

type WSEvent = {
  type: string;
  payload?: Record<string, unknown>;
  project_id?: string;
};

export function useWebSocket() {
  const { isAuthenticated, logout } = useAuth();
  const queryClient = useQueryClient();
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const reconnectDelay = useRef(1000);
  const connectRef = useRef<() => void>(null);

  const handleEvent = useCallback(
    (evt: WSEvent) => {
      switch (evt.type) {
        case 'force_logout':
          toast.error('Session expired', { description: 'Please log in again.' });
          logout();
          break;
        case 'task_created':
        case 'task_updated':
        case 'task_deleted': {
          const projectId = evt.project_id ?? (evt.payload?.project_id as string);
          if (projectId) {
            queryClient.invalidateQueries({ queryKey: ['tasks', projectId] });
            queryClient.invalidateQueries({ queryKey: ['project', projectId] });
          }
          queryClient.invalidateQueries({ queryKey: ['dashboard'] });
          queryClient.invalidateQueries({ queryKey: ['workspace-stats'] });
          queryClient.invalidateQueries({ queryKey: ['projects'] });
          queryClient.invalidateQueries({ queryKey: ['my-stats'] });
          queryClient.invalidateQueries({ queryKey: ['org-dashboard'] });
          break;
        }
      }
    },
    [queryClient, logout],
  );

  const connect = useCallback(() => {
    const token = localStorage.getItem('token');
    if (!token) return;

    let wsUrl: string;
    if (WS_BASE_URL) {
      wsUrl = `${WS_BASE_URL}/ws?token=${token}`;
    } else {
      const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      wsUrl = `${proto}//${window.location.host}/ws?token=${token}`;
    }

    const socket = new WebSocket(wsUrl);
    wsRef.current = socket;

    socket.onopen = () => {
      reconnectDelay.current = 1000;
    };

    socket.onmessage = (e) => {
      try {
        const data: WSEvent = JSON.parse(e.data);
        handleEvent(data);
      } catch {
        // ignore malformed messages
      }
    };

    socket.onclose = (e) => {
      wsRef.current = null;
      if (e.code === 1000) return;
      reconnectTimer.current = setTimeout(() => {
        reconnectDelay.current = Math.min(reconnectDelay.current * 2, 30000);
        connectRef.current?.();
      }, reconnectDelay.current);
    };

    socket.onerror = () => {
      socket.close();
    };
  }, [handleEvent]);

  useEffect(() => {
    connectRef.current = connect;
  }, [connect]);

  useEffect(() => {
    if (isAuthenticated) {
      connect();
    }
    return () => {
      if (reconnectTimer.current) clearTimeout(reconnectTimer.current);
      wsRef.current?.close(1000, 'component unmount');
    };
  }, [isAuthenticated, connect]);
}
