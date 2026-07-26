import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { Fragancia, FraganciasFilters } from '../types';

export function useFragancias(filters: FraganciasFilters) {
  return useQuery({
    queryKey: ['fragancias', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        q: filters.q,
        genero: filters.genero === 'all' ? '' : filters.genero,
        activo: filters.activo === 'all' ? 'all' : filters.activo,
        stock_bajo: filters.stock_bajo,
      });
      const res = await api.get<ApiListEnvelope<Fragancia>>('/fragancias', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
