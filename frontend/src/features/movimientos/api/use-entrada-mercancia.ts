import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { TipoItem, Ubicacion } from '@/types/domain';
import { invalidateInventory } from './invalidate';

export interface EntradaItemPayload {
  tipo_item: TipoItem;
  item_id: number;
  ubicacion: Ubicacion;
  cantidad: string;
  motivo?: string;
}

export function useEntradaMercancia() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (items: EntradaItemPayload[]) => {
      const res = await api.post('/movimientos/entrada-mercancia', { items });
      return res.data.data;
    },
    onSuccess: () => invalidateInventory(queryClient),
  });
}
