import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { MetodoPago } from '../types';

export function useMetodoPago(id: number) {
  return useQuery({
    queryKey: ['metodos-pago', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<MetodoPago>>(`/metodos-pago/${id}`);
      return res.data.data;
    },
    enabled: Number.isFinite(id) && id > 0,
  });
}
