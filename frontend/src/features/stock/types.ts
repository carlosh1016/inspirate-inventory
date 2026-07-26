import type { TipoItem } from '@/types/domain';

// Mirrors handlers/stock/dto.go StockItemResponse.
export interface StockItem {
  tipo_item: TipoItem;
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

export type StockFilters = {
  page: number;
  tipo_item: string;
  ubicacion: string;
  stock_bajo: boolean;
  stock_cero: boolean;
};
