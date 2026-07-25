import { useQuery } from '@tanstack/react-query';

import { api } from '@/lib/api';
import type { ApiListEnvelope } from '@/types/api';

// /stock/alertas is paginated (WriteList). The KPI count uses meta.total, not
// the page length. Shape verified against handlers/stock/dto.go StockItemResponse.
export interface StockAlerta {
  tipo_item: string;
  item_id: number;
  nombre: string;
  detalle_extra: string;
  stock_vitrina: string;
  stock_bodega: string;
  stock_total: string;
  minimo: string;
  bajo_minimo: boolean;
  unidad: string;
}

export function useAlertasStock() {
  return useQuery({
    queryKey: ['stock', 'alertas'],
    queryFn: async () => {
      const res = await api.get<ApiListEnvelope<StockAlerta>>('/stock/alertas?page=1&page_size=50');
      return { items: res.data.data, total: res.data.meta.total };
    },
  });
}
