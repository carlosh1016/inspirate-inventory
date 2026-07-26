import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Producto } from '../types';

export function useProducto(id: number) {
  return useQuery({
    queryKey: ['productos', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<Producto>>(`/productos/${id}`);
      return res.data.data;
    },
    enabled: Number.isFinite(id) && id > 0,
  });
}
