import { formatGramos } from '@/lib/formatters';
import { cn } from '@/lib/utils';

interface Props {
  vitrina: string;
  bodega: string;
  total: string;
  minimo?: string;
  /** "gramos" formats with formatGramos; otherwise shows the raw number. */
  unidad?: string;
  className?: string;
}

function fmt(value: string, unidad: string): string {
  if (unidad === 'gramos') return formatGramos(value);
  const n = Number.parseFloat(value);
  return Number.isNaN(n) ? value : String(n);
}

// Vitrina / Bodega / Total tiles for an item's current stock.
export function StockSummaryCard({ vitrina, bodega, total, minimo, unidad = 'gramos', className }: Props) {
  const tiles = [
    { label: 'Vitrina', value: vitrina },
    { label: 'Bodega', value: bodega },
    { label: 'Total', value: total },
  ];

  return (
    <div className={cn('space-y-2', className)}>
      <div className="grid grid-cols-3 gap-3">
        {tiles.map((tile) => (
          <div key={tile.label} className="rounded-lg border border-border p-3">
            <p className="text-xs text-muted-foreground">{tile.label}</p>
            <p className="mt-1 text-lg font-semibold tabular-nums">{fmt(tile.value, unidad)}</p>
          </div>
        ))}
      </div>
      {minimo !== undefined && (
        <p className="text-sm text-muted-foreground">Mínimo: {fmt(minimo, unidad)}</p>
      )}
    </div>
  );
}
