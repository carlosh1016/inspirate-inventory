import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';

// Shape verified against handlers/movimientos/dto.go MovimientoResponse. `item`
// and `usuario` are nested objects; quantities are decimal strings.
export interface Movimiento {
  id: number;
  tipo_item: string;
  item: { id: number; nombre: string };
  tipo: string;
  ubicacion: string;
  cantidad: string;
  stock_anterior: string;
  stock_posterior: string;
  usuario: { id: number; nombre_completo?: string };
  venta_id?: number;
  motivo?: string;
  created_at: string;
}

export function useUltimosMovimientos() {
  return useQuery({
    queryKey: ['movimientos', 'ultimos'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<Movimiento>>('/movimientos?page=1&page_size=10');
      return res.data.data;
    },
  });
}
