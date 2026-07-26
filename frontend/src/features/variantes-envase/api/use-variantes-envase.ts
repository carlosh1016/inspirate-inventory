import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { VarianteEnvase, VariantesEnvaseFilters } from '../types';

export function useVariantesEnvase(filters: VariantesEnvaseFilters) {
  return useQuery({
    queryKey: ['variantes-envase', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        q: filters.q,
        modelo_envase_id: filters.modelo_envase_id > 0 ? filters.modelo_envase_id : '',
        activo: filters.activo === 'all' ? 'all' : filters.activo,
        stock_bajo: filters.stock_bajo,
      });
      const res = await api.get<ApiListEnvelope<VarianteEnvase>>('/variantes-envase', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
