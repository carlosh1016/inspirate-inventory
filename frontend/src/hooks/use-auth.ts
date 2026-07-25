'use client';

import { useAuthStore } from '@/stores/auth-store';

/** Convenience selector over the auth store. */
export function useAuth() {
  const usuario = useAuthStore((s) => s.usuario);
  const accessToken = useAuthStore((s) => s.accessToken);
  const isInitializing = useAuthStore((s) => s.isInitializing);
  return {
    usuario,
    accessToken,
    isInitializing,
    isAuthenticated: usuario !== null,
  };
}
