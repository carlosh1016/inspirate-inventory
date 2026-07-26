import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import type { ModeloEnvase } from '../types';

// id -> full ModeloEnvase map, for attaching prices/equiv_gramos to a variante
// in the venta form. Modelos are few, so one page suffices.
export function useModelosFullLookup() {
  return useQuery({
    queryKey: ['modelos-envase', 'full-lookup'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<ModeloEnvase>>('/modelos-envase', {
        params: { activo: 'all', page_size: 100 },
      });
      const map = new Map<number, ModeloEnvase>();
      res.data.data.forEach((m) => map.set(m.id, m));
      return map;
    },
  });
}
