import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';

// Shape verified against handlers/ventas/dto.go ResumenHoyResponse. Monetary
// values are decimal strings.
export interface PorMetodoPago {
  efectivo: string;
  nequi: string;
  daviplata: string;
  transferencia: string;
  otros: string;
}

export interface VentaPorVendedora {
  usuario_id: number;
  nombre_completo: string;
  ventas_count: number;
  total: string;
}

export interface TopFragancia {
  id: number;
  nombre_comercial: string;
  gramos_vendidos: string;
  monto_vendido: string;
}

export interface ResumenHoy {
  fecha: string;
  ventas_count: number;
  total_dia: string;
  descuento_total: string;
  por_metodo_pago: PorMetodoPago;
  por_vendedora: VentaPorVendedora[];
  top_fragancias: TopFragancia[];
}

export function useResumenHoy() {
  return useQuery({
    queryKey: ['ventas', 'resumen-hoy'],
    queryFn: async () => {
      const res = await api.get<ApiEnvelope<ResumenHoy>>('/ventas/hoy/resumen');
      return res.data.data;
    },
    refetchInterval: 60_000,
  });
}
