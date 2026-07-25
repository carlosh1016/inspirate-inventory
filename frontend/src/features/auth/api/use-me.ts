import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { UsuarioApi } from '@/types/domain';

export function useMe(enabled = true) {
  return useQuery({
    queryKey: ['auth', 'me'],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<UsuarioApi>>('/auth/me');
      return res.data.data;
    },
    enabled,
  });
}
