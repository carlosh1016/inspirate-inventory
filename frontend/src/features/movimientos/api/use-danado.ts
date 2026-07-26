import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { TipoItem, Ubicacion } from '@/types/domain';
import { invalidateInventory } from './invalidate';

export interface DanadoPayload {
  tipo_item: TipoItem;
  item_id: number;
  ubicacion: Ubicacion;
  cantidad: string;
  motivo: string;
}

export function useDanado() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload: DanadoPayload) => {
      const res = await api.post('/movimientos/danado', payload);
      return res.data.data;
    },
    onSuccess: () => invalidateInventory(queryClient),
  });
}
