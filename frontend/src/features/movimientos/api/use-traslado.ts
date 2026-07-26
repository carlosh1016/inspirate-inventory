import { useMutation, useQueryClient } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { TipoItem } from '@/types/domain';
import { invalidateInventory } from './invalidate';

export interface TrasladoItemPayload {
  tipo_item: TipoItem;
  item_id: number;
  cantidad: string;
}

export function useTraslado() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (items: TrasladoItemPayload[]) => {
      const res = await api.post('/movimientos/traslado', { items });
      return res.data.data;
    },
    onSuccess: () => invalidateInventory(queryClient),
  });
}
