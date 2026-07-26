import { formatCurrency, formatGramos } from '@/lib/formatters';
import type { TipoLinea, VentaItemDetalle } from '../types';

const TIPO_LABEL: Record<TipoLinea, string> = {
  envase_con_fragancia: 'Envase con fragancia',
  recarga: 'Recarga',
  envase_solo: 'Envase solo',
  producto_otro: 'Otro',
};

function itemTitulo(item: VentaItemDetalle): string {
  switch (item.tipo_linea) {
    case 'envase_con_fragancia':
    case 'recarga':
      return [item.variante_envase_nombre, item.fragancia_nombre].filter(Boolean).join(' + ');
    case 'envase_solo':
      return item.variante_envase_nombre ?? 'Envase';
    case 'producto_otro':
      return item.producto_nombre ?? 'Producto';
  }
}

function itemDetalleLinea(item: VentaItemDetalle): string {
  if (item.tipo_linea === 'envase_con_fragancia' || item.tipo_linea === 'recarga') {
    const gramos = item.gramos_fragancia ? formatGramos(item.gramos_fragancia) : '';
    return `${gramos} × ${item.cantidad}`;
  }
  return `× ${item.cantidad}`;
}

export function VentaItemsList({ items }: { items: VentaItemDetalle[] }) {
  return (
    <div className="space-y-3">
      {items.map((item) => (
        <div key={item.id} className="rounded-lg border border-border p-4">
          <div className="flex items-start justify-between gap-4">
            <div className="min-w-0 space-y-1">
              <p className="text-xs text-muted-foreground">{TIPO_LABEL[item.tipo_linea]}</p>
              <p className="font-medium">{itemTitulo(item)}</p>
              <p className="text-sm text-muted-foreground">{itemDetalleLinea(item)}</p>
              {item.feromona_nombre && (
                <p className="text-sm text-muted-foreground">+ Feromonas ({item.feromona_nombre})</p>
              )}
            </div>
            <p className="shrink-0 font-medium tabular-nums">{formatCurrency(item.subtotal)}</p>
          </div>
        </div>
      ))}
    </div>
  );
}
