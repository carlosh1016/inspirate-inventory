import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import type { MetodoPago } from '../types';

export function useMetodosPago() {
  return useQuery({
    queryKey: ['metodos-pago'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<MetodoPago>>('/metodos-pago', {
        params: { page_size: 100 },
      });
      return res.data.data;
    },
  });
}
