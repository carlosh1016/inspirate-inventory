import { formatCurrency } from '@/lib/formatters';
import type { CalculoVenta } from '../hooks/use-calculo-venta';

export function VentaTotales({ calculo }: { calculo: CalculoVenta }) {
  const { subtotal, descuentoPct, descuentoMonto, total } = calculo;
  return (
    <div className="space-y-1.5 text-sm">
      <div className="flex justify-between">
        <span className="text-muted-foreground">Subtotal</span>
        <span className="tabular-nums">{formatCurrency(subtotal)}</span>
      </div>
      {descuentoPct > 0 && (
        <div className="flex justify-between text-success">
          <span>
            Descuento ({descuentoPct}%){' '}
            <span className="text-xs">(por superar {formatCurrency(descuentoPct === 7 ? 100000 : 50000)})</span>
          </span>
          <span className="tabular-nums">-{formatCurrency(descuentoMonto)}</span>
        </div>
      )}
      <div className="mt-2 flex items-baseline justify-between border-t border-border pt-2">
        <span className="font-medium">Total</span>
        <span className="text-2xl font-semibold tabular-nums">{formatCurrency(total)}</span>
      </div>
    </div>
  );
}
