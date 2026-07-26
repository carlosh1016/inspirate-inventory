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
    rojo: 'bg-destructive/10 text-destructive',
    ambar: 'bg-warning/10 text-warning',
    verde: 'bg-success/10 text-success',
  } as const;

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 rounded-md px-2 py-0.5 text-xs font-medium whitespace-nowrap',
        styles[level],
        className,
      )}
    >
      {level === 'rojo' && <AlertTriangle className="size-3" />}
      {label}
    </span>
  );
}
