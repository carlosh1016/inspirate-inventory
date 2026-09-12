import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import type { ConteoListItem } from '../types';

// GET /conteos-inventario — admin only.
export function useConteos(page: number) {
  return useQuery({
    queryKey: ['conteos-inventario', 'list', page],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<ConteoListItem>>('/conteos-inventario', {
        params: { page },
      });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
