import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { VarianteEnvase } from '../types';

export function useVarianteEnvase(id: number) {
  return useQuery({
    queryKey: ['variantes-envase', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<VarianteEnvase>>(`/variantes-envase/${id}`);
      return res.data.data;
    },
    enabled: Number.isFinite(id) && id > 0,
  });
}
