import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { TipoItem, Ubicacion } from '@/types/domain';
import { invalidateInventory } from './invalidate';
import type { Movimiento } from '../types';

export interface AjustePayload {
  tipo_item: TipoItem;
  item_id: number;
  ubicacion: Ubicacion;
  cantidad_nueva: string;
  motivo: string;
}

// Ajuste returns 201 with a Movimiento, or 200 { mensaje } when the stock
// already matched the requested quantity (no movimiento created).
export interface AjusteResult {
  movimiento: Movimiento | null;
  mensaje?: string;
}

export function useAjuste() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload: AjustePayload): Promise<AjusteResult> => {
      const res = await api.post('/movimientos/ajuste', payload);
      const data = res.data.data;
      if (data && typeof data === 'object' && 'mensaje' in data) {
        return { movimiento: null, mensaje: data.mensaje as string };
      }
      return { movimiento: data as Movimiento };
    },
    onSuccess: () => invalidateInventory(queryClient),
  });
}
