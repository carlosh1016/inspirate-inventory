'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

import { useAuthStore } from '@/stores/auth-store';

// Renders children only for an authenticated user; otherwise redirects to
// /login once initialization has settled.
export function RequireAuth({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const usuario = useAuthStore((s) => s.usuario);
  const isInitializing = useAuthStore((s) => s.isInitializing);

  useEffect(() => {
    if (!isInitializing && !usuario) {
      router.replace('/login');
    }
  }, [isInitializing, usuario, router]);

  if (!usuario) return null;
  return <>{children}</>;
}
