'use client';

import { useAuthStore } from '@/stores/auth-store';
import type { Rol } from '@/types/domain';

/** Role helpers derived from the current session. */
export function usePermission() {
  const usuario = useAuthStore((s) => s.usuario);
  return {
    rol: usuario?.rol,
    isAdmin: usuario?.rol === 'admin',
    isVendedora: usuario?.rol === 'vendedora',
    hasRole: (role: Rol) => usuario?.rol === role,
  };
}
