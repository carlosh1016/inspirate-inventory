import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Cuadre } from '../types';

// GET /cuadres-caja/hoy returns { data: null } when no cuadre exists today.
export function useCuadreHoy() {
  return useQuery({
    queryKey: ['cuadres-caja', 'hoy'],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<Cuadre | null>>('/cuadres-caja/hoy');
      return res.data.data;
    },
  });
}
