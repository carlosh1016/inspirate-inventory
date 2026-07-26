import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';

interface UsuarioRow {
  id: number;
  nombre_completo: string;
}

// Lightweight usuarios list for the movimientos "usuario" filter. /usuarios is
// admin-only, so this is gated by `enabled`.
export function useUsuariosSelect(enabled: boolean) {
  return useQuery({
    queryKey: ['usuarios', 'select'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<UsuarioRow>>('/usuarios', {
        params: { page_size: 100 },
      });
      return res.data.data;
    },
    enabled,
  });
}
