import { parseApiError } from '@/lib/errors';

// The "stock insuficiente" error carries `extra` as a flat array (see backend
// domain/stock/errors.go StockInsuficienteItem). Match rows to form items by
// (tipo_item, item_id).
export interface StockInsuficienteItem {
  tipo_item: string;
  item_id: number;
  nombre: string;
  requerido: string;
  disponible: string;
}

export function parseStockInsuficiente(error: unknown): StockInsuficienteItem[] | null {
  const { extra } = parseApiError(error);
  if (!Array.isArray(extra)) return null;
  return extra as StockInsuficienteItem[];
}
