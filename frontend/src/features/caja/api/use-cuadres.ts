import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { CuadreListItem, CuadresFilters } from '../types';

// GET /cuadres-caja — admin only. Dates are YYYY-MM-DD (backend parseOptionalDate).
export function useCuadres(filters: CuadresFilters) {
  return useQuery({
    queryKey: ['cuadres-caja', 'list', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        estado: filters.estado === 'all' ? '' : filters.estado,
        fecha_desde: filters.fecha_desde,
        fecha_hasta: filters.fecha_hasta,
      });
      const res = await api.get<ApiListEnvelope<CuadreListItem>>('/cuadres-caja', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
