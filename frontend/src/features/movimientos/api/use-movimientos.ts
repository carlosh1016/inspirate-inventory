import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import { buildQueryParams, dateToIsoEnd, dateToIsoStart } from '@/lib/query-params';
import type { ApiListEnvelope } from '@/types/api';
import type { Movimiento } from '../types';

export type MovimientosFilters = {
  page: number;
  tipo: string;
  tipo_item: string;
  ubicacion: string;
  item_id: number;
  item_nombre: string;
  usuario_id: number;
  fecha_desde: string;
  fecha_hasta: string;
};

export function useMovimientos(filters: MovimientosFilters) {
  return useQuery({
    queryKey: ['movimientos', 'list', filters],
    queryFn: async () => {
      const params = buildQueryParams({
        page: filters.page,
        tipo: filters.tipo === 'all' ? '' : filters.tipo,
        tipo_item: filters.tipo_item === 'all' ? '' : filters.tipo_item,
        ubicacion: filters.ubicacion === 'all' ? '' : filters.ubicacion,
        item_id: filters.item_id > 0 ? filters.item_id : '',
        usuario_id: filters.usuario_id > 0 ? filters.usuario_id : '',
        fecha_desde: dateToIsoStart(filters.fecha_desde) ?? '',
        fecha_hasta: dateToIsoEnd(filters.fecha_hasta) ?? '',
      });
      const res = await api.get<ApiListEnvelope<Movimiento>>('/movimientos', { params });
      return { items: res.data.data, meta: res.data.meta };
    },
  });
}
