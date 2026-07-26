import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { Producto, ProductosFilters } from '../types';

export function useProductos(filters: ProductosFilters) {
  return useQuery({
    queryKey: ['productos', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        q: filters.q,
        categoria: filters.categoria === 'all' ? '' : filters.categoria,
        activo: filters.activo === 'all' ? 'all' : filters.activo,
        stock_bajo: filters.stock_bajo,
      });
      const res = await api.get<ApiListEnvelope<Producto>>('/productos', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
