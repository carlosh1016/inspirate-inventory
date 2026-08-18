import { AlertTriangle } from 'lucide-react';

import { formatGramos } from '@/lib/formatters';
import { cn } from '@/lib/utils';

interface Props {
  total: string | number;
  minimo: string | number;
  /** "gramos" formats with formatGramos; anything else shows "N unidades". */
  unidad?: string;
  className?: string;
}

function toNumber(value: string | number): number {
  const n = typeof value === 'string' ? Number.parseFloat(value) : value;
  return Number.isNaN(n) ? 0 : n;
}

// Visual stock level: rojo (<= mínimo, con icono), ámbar (<= 2x mínimo),
// verde (por encima).
export function StockBadge({ total, minimo, unidad = 'gramos', className }: Props) {
  const t = toNumber(total);
  const m = toNumber(minimo);
  const level = t <= m ? 'rojo' : t <= m * 2 ? 'ambar' : 'verde';

  const label =
    unidad === 'gramos' ? formatGramos(t) : `${t} ${t === 1 ? 'unidad' : 'unidades'}`;

  const styles = {
    rojo: 'border-destructive/20 bg-destructive/10 text-destructive',
    ambar: 'border-warning/20 bg-warning/10 text-warning',
    verde: 'border-success/20 bg-success/10 text-success',
  } as const;

  const dotStyles = {
    rojo: 'bg-destructive',
    ambar: 'bg-warning',
    verde: 'bg-success',
  } as const;

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 font-mono text-xs font-semibold whitespace-nowrap',
        styles[level],
        className,
      )}
    >
      {level === 'rojo' ? (
        <AlertTriangle className="size-3" />
      ) : (
        <span className={cn('size-1.5 shrink-0 rounded-full', dotStyles[level])} />
      )}
      {label}
    </span>
  );
}
