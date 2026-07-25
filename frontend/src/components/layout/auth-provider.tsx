'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/auth-store';
import { toUsuarioSession, type UsuarioApi } from '@/types/domain';

// Restores the session on mount: the access token lives only in memory, so on a
// fresh load we POST /auth/refresh (the HttpOnly cookie is sent because its path
// matches) and, on success, GET /auth/me. On failure we clear and go to /login.
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const isInitializing = useAuthStore((s) => s.isInitializing);
  const accessToken = useAuthStore((s) => s.accessToken);

  useEffect(() => {
    let cancelled = false;
    const store = useAuthStore.getState();

    async function initSession() {
      try {
        const refreshRes = await api.post('/auth/refresh');
        const token: string | undefined = refreshRes.data?.data?.access_token;
        if (!token) throw new Error('no token');
        store.setAccessToken(token);

        const meRes = await api.get('/auth/me');
        const usuario = meRes.data?.data as UsuarioApi;
        if (!cancelled) store.setSession(token, toUsuarioSession(usuario));
      } catch {
        if (!cancelled) {
          store.clearSession();
          router.replace('/login');
        }
      } finally {
        if (!cancelled) store.setInitializing(false);
      }
    }

    if (accessToken) {
      store.setInitializing(false);
    } else {
      void initSession();
    }

    return () => {
      cancelled = true;
    };
    // Runs once on mount; the session is (re)validated here, not on every render.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (isInitializing) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-sm text-muted-foreground">Cargando...</p>
      </div>
    );
  }

  return <>{children}</>;
}
