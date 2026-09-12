import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Conteo } from '../types';

// GET /conteos-inventario/:id — admin only.
export function useConteo(id: number) {
  return useQuery({
    queryKey: ['conteos-inventario', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<Conteo>>(`/conteos-inventario/${id}`);
      return res.data.data;
    },
    enabled: Number.isFinite(id) && id > 0,
  });
}
