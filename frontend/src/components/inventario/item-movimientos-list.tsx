'use client';

import { useItemMovimientos } from '@/features/movimientos/api/use-item-movimientos';
import { TIPO_MOVIMIENTO_META } from '@/features/movimientos/tipo-meta';
import { formatGramos, formatRelative } from '@/lib/formatters';
import { cn } from '@/lib/utils';
import type { TipoItem } from '@/types/domain';

interface Props {
  tipoItem: TipoItem;
  itemId: number;
  unidad?: string;
}

function fmtCantidad(cantidad: string, unidad: string): string {
  const n = Number.parseFloat(cantidad);
  const sign = n > 0 ? '+' : ''; // negatives already carry '-'
  const body = unidad === 'gramos' ? formatGramos(n) : String(n);
  return `${sign}${body}`;
}

// Recent movimientos for one catalog item (detail pages).
export function ItemMovimientosList({ tipoItem, itemId, unidad = 'gramos' }: Props) {
  const { data, isLoading } = useItemMovimientos(tipoItem, itemId);

  if (isLoading) {
    return <p className="text-sm text-muted-foreground">Cargando movimientos…</p>;
  }
  if (!data || data.length === 0) {
    return <p className="text-sm text-muted-foreground">Sin movimientos registrados.</p>;
  }

  return (
    <ul className="divide-y divide-border">
      {data.map((mov) => {
        const meta = TIPO_MOVIMIENTO_META[mov.tipo];
        const negative = Number.parseFloat(mov.cantidad) < 0;
        return (
          <li key={mov.id} className="flex items-center justify-between gap-3 py-2 text-sm">
            <div className="min-w-0">
              <span className={cn('rounded-md px-1.5 py-0.5 text-xs font-medium', meta.className)}>
                {meta.label}
              </span>
              <span className="ml-2 text-muted-foreground">
                {formatRelative(mov.created_at)}
                {mov.usuario.nombre_completo ? ` · ${mov.usuario.nombre_completo}` : ''}
              </span>
            </div>
            <span
              className={cn('shrink-0 tabular-nums', negative ? 'text-destructive' : 'text-success')}
            >
              {fmtCantidad(mov.cantidad, unidad)}
            </span>
          </li>
        );
      })}
    </ul>
  );
}
