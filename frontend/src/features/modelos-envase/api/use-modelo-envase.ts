import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { ModeloEnvase } from '../types';

export function useModeloEnvase(id: number) {
  return useQuery({
    queryKey: ['modelos-envase', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<ModeloEnvase>>(`/modelos-envase/${id}`);
      return res.data.data;
    },
    enabled: Number.isFinite(id) && id > 0,
  });
}
