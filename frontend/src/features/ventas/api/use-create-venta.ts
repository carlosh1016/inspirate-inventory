import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiEnvelope } from '@/types/api';
import type { TipoLinea, VentaDetallada } from '../types';

export interface CreateVentaItemPayload {
  tipo_linea: TipoLinea;
  fragancia_id?: number;
  variante_envase_id?: number;
  producto_id?: number;
  feromona_producto_id?: number;
  gramos_fragancia?: string;
  cantidad: number;
}

export interface CreateVentaPayload {
  metodo_pago_id: number;
  observaciones?: string | null;
  items: CreateVentaItemPayload[];
}

export function useCreateVenta() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({
      payload,
      idempotencyKey,
    }: {
      payload: CreateVentaPayload;
      idempotencyKey: string;
    }) => {
      const res = await api.post<ApiEnvelope<VentaDetallada>>('/ventas', payload, {
        headers: { 'X-Idempotency-Key': idempotencyKey },
      });
      return res.data.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ventas'] });
      queryClient.invalidateQueries({ queryKey: ['stock'] });
      queryClient.invalidateQueries({ queryKey: ['cuadres-caja'] });
      queryClient.invalidateQueries({ queryKey: ['movimientos'] });
    },
  });
}
