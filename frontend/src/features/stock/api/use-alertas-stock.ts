import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { StockItem } from '../types';

// Paginated list of items below their minimum (GET /stock/alertas).
export function useAlertasStock(page: number) {
  return useQuery({
    queryKey: ['stock', 'alertas', page],
    queryFn: async () => {
      const params = buildQueryParams({ page });
      const res = await api.get<ApiListEnvelope<StockItem>>('/stock/alertas', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
