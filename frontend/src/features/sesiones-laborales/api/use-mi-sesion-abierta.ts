import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import type { Sesion } from '../types';

// The vendedora's currently-open sesion (backend auto-scopes to self). Returns
// null when there's no open turno.
export function useMiSesionAbierta() {
  return useQuery({
    queryKey: ['sesiones-laborales', 'abierta'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<Sesion>>('/sesiones-laborales', {
        params: { abiertas: true, page_size: 1 },
      });
      return res.data.data[0] ?? null;
    },
  });
}
