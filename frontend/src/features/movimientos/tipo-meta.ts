import type { TipoItem } from '@/types/domain';
import type { TipoMovimiento } from './types';

// Display label + badge color per movimiento type. Colors follow the semantic
// tokens: entradas verde, salidas/dañado rojo, ajuste/corrección azul, etc.
export const TIPO_MOVIMIENTO_META: Record<
  TipoMovimiento,
  { label: string; className: string }
> = {
  entrada_mercancia: { label: 'Entrada', className: 'bg-success/10 text-success' },
  traslado_bodega_vitrina: { label: 'Traslado', className: 'bg-info/10 text-info' },
  venta: { label: 'Venta', className: 'bg-destructive/10 text-destructive' },
  ajuste: { label: 'Ajuste', className: 'bg-info/10 text-info' },
  danado: { label: 'Dañado', className: 'bg-destructive/10 text-destructive' },
  devolucion: { label: 'Devolución', className: 'bg-success/10 text-success' },
  correccion: { label: 'Corrección', className: 'bg-info/10 text-info' },
};

export const TIPO_ITEM_LABEL: Record<TipoItem, string> = {
  fragancia: 'Fragancia',
  variante_envase: 'Variante de envase',
  producto: 'Producto',
};
