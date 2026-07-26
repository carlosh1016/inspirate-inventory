import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import { modeloEnvaseLabel } from './search-modelos-envase';
import type { ModeloEnvase } from '../types';

// id -> display label map for showing a modelo's name where only its id is
// available (e.g. the variantes table). Modelos are few, so one page suffices.
export function useModelosLookup() {
  return useQuery({
    queryKey: ['modelos-envase', 'lookup'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<ModeloEnvase>>('/modelos-envase', {
        params: { activo: 'all', page_size: 100 },
      });
      const map = new Map<number, string>();
      res.data.data.forEach((m) => map.set(m.id, modeloEnvaseLabel(m)));
      return map;
    },
  });
}
