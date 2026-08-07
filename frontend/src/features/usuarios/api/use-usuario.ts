import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { UsuarioApi } from '@/types/domain';

export function useUsuario(id: number) {
  return useQuery({
    queryKey: ['usuarios', 'detail', id],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<UsuarioApi>>(`/usuarios/${id}`);
      return res.data.data;
    },
    enabled: id > 0,
  });
}
