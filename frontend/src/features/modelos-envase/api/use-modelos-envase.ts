import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { ModeloEnvase, ModelosEnvaseFilters } from '../types';

export function useModelosEnvase(filters: ModelosEnvaseFilters) {
  return useQuery({
    queryKey: ['modelos-envase', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        q: filters.q,
        activo: filters.activo === 'all' ? 'all' : filters.activo,
      });
      const res = await api.get<ApiListEnvelope<ModeloEnvase>>('/modelos-envase', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
