import type { TipoItem, Ubicacion } from '@/types/domain';

export type TipoMovimiento =
  | 'entrada_mercancia'
  | 'traslado_bodega_vitrina'
  | 'venta'
  | 'ajuste'
  | 'danado'
  | 'devolucion'
  | 'correccion';

export interface MovimientoItemRef {
  id: number;
  nombre: string;
}

export interface MovimientoUsuarioRef {
  id: number;
  nombre_completo?: string;
}

// Mirrors handlers/movimientos/dto.go MovimientoResponse.
export interface Movimiento {
  id: number;
  tipo_item: TipoItem;
  item: MovimientoItemRef;
  tipo: TipoMovimiento;
  ubicacion: Ubicacion;
  cantidad: string;
  stock_anterior: string;
  stock_posterior: string;
  usuario: MovimientoUsuarioRef;
  venta_id?: number | null;
  motivo?: string | null;
  created_at: string;
}
