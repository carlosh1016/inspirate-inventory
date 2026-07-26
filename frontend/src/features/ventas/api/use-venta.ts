import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { VentaDetallada } from '../types';

export function useVenta(id: number) {
  return useQuery({
    queryKey: ['ventas', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<VentaDetallada>>(`/ventas/${id}`);
      return res.data.data;
    },
    enabled: Number.isFinite(id) && id > 0,
  });
}
