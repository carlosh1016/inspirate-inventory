import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Cuadre } from '../types';

// GET /cuadres-caja/:id — admin only.
export function useCuadre(id: number) {
  return useQuery({
    queryKey: ['cuadres-caja', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<Cuadre>>(`/cuadres-caja/${id}`);
      return res.data.data;
    },
    enabled: Number.isFinite(id) && id > 0,
  });
}
