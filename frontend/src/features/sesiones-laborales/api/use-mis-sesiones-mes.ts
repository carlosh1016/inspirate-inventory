import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { Sesion } from '../types';

// The vendedora's sesiones for the current month (backend auto-scopes to self).
export function useMisSesionesMes() {
  const now = new Date();
  const inicioMes = new Date(now.getFullYear(), now.getMonth(), 1);
  const finMes = new Date(now.getFullYear(), now.getMonth() + 1, 0, 23, 59, 59);

  return useQuery({
    queryKey: ['sesiones-laborales', 'mes', inicioMes.toISOString().slice(0, 7)],
    queryFn: async () => {
      const params = buildQueryParams({
        fecha_desde: inicioMes.toISOString(),
        fecha_hasta: finMes.toISOString(),
        page_size: 100,
      });
      const res = await api.get<ApiListEnvelope<Sesion>>('/sesiones-laborales', { params });
      return res.data.data;
    },
  });
}
