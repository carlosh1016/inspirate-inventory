import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { StockFilters, StockItem } from '../types';

export function useStock(filters: StockFilters) {
  return useQuery({
    queryKey: ['stock', 'list', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        tipo_item: filters.tipo_item === 'all' ? '' : filters.tipo_item,
        ubicacion: filters.ubicacion === 'all' ? '' : filters.ubicacion,
        stock_bajo: filters.stock_bajo,
        stock_cero: filters.stock_cero,
      });
      const res = await api.get<ApiListEnvelope<StockItem>>('/stock', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
