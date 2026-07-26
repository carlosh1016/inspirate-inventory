import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams, dateToIsoEnd, dateToIsoStart } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { VentaListItem, VentasFilters } from '../types';

export function useVentas(filters: VentasFilters) {
  return useQuery({
    queryKey: ['ventas', 'list', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        metodo_pago_id: filters.metodo_pago_id > 0 ? filters.metodo_pago_id : '',
        usuario_id: filters.usuario_id > 0 ? filters.usuario_id : '',
        fecha_desde: dateToIsoStart(filters.fecha_desde) ?? '',
        fecha_hasta: dateToIsoEnd(filters.fecha_hasta) ?? '',
        con_descuento: filters.con_descuento,
        total_min: filters.total_min,
        total_max: filters.total_max,
      });
      const res = await api.get<ApiListEnvelope<VentaListItem>>('/ventas', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
