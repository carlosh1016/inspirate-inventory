import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { UsuarioApi } from '@/types/domain';
import type { ApiListEnvelope } from '@/types/api';
import type { UsuariosFilters } from '../types';

export function useUsuarios(filters: UsuariosFilters) {
  return useQuery({
    queryKey: ['usuarios', 'list', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        rol: filters.rol === 'all' ? '' : filters.rol,
        activo: filters.activo === 'all' ? '' : filters.activo,
      });
      const res = await api.get<ApiListEnvelope<UsuarioApi>>('/usuarios', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
