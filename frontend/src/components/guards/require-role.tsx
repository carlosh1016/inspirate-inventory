'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

import { useAuthStore } from '@/stores/auth-store';
import type { Rol } from '@/types/domain';

interface Props {
  role: Rol;
  children: React.ReactNode;
}

// Renders children only when the session role matches; a mismatched role is
// sent to /acceso-denegado.
export function RequireRole({ role, children }: Props) {
  const router = useRouter();
  const usuario = useAuthStore((s) => s.usuario);

  useEffect(() => {
    if (usuario && usuario.rol !== role) {
      router.replace('/acceso-denegado');
    }
  }, [usuario, role, router]);

  if (!usuario || usuario.rol !== role) return null;
  return <>{children}</>;
}
