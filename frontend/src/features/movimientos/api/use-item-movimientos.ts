import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';
import type { TipoItem } from '@/types/domain';
import type { Movimiento } from '../types';

// Recent movimientos for a single catalog item, used in item detail pages.
export function useItemMovimientos(tipoItem: TipoItem, itemId: number, pageSize = 5) {
  return useQuery({
    queryKey: ['movimientos', 'item', tipoItem, itemId, pageSize],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<Movimiento>>('/movimientos', {
        params: { tipo_item: tipoItem, item_id: itemId, page_size: pageSize },
      });
      return res.data.data;
    },
    enabled: Number.isFinite(itemId) && itemId > 0,
  });
}
