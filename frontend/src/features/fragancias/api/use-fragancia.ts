import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { Fragancia } from '../types';

export function useFragancia(id: number) {
  return useQuery({
    queryKey: ['fragancias', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<Fragancia>>(`/fragancias/${id}`);
      return res.data.data;
    },
    enabled: Number.isFinite(id) && id > 0,
  });
}
