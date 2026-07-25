import { EmptyState } from '@/components/feedback/empty-state';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import type { Movimiento } from '@/features/dashboard/api/use-ultimos-movimientos';
import { formatRelative } from '@/lib/formatters';

const tipoLabels: Record<string, string> = {
  entrada_mercancia: 'Entrada',
  traslado_bodega_vitrina: 'Traslado',
  venta: 'Venta',
  ajuste: 'Ajuste',
  danado: 'Dañado',
  devolucion: 'Devolución',
  correccion: 'Corrección',
};

function tipoLabel(tipo: string): string {
  return tipoLabels[tipo] ?? tipo;
}

interface Props {
  movimientos: Movimiento[] | undefined;
  isLoading: boolean;
}

export function UltimosMovimientosCard({ movimientos, isLoading }: Props) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Últimos movimientos</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="space-y-2">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-10" />
            ))}
          </div>
        ) : !movimientos || movimientos.length === 0 ? (
          <EmptyState title="Sin movimientos recientes" />
        ) : (
          <ul className="divide-y divide-border">
            {movimientos.map((m) => (
              <li key={m.id} className="flex items-center justify-between gap-3 py-2">
                <div className="min-w-0">
                  <p className="truncate text-sm font-medium">{m.item.nombre}</p>
                  <p className="text-xs text-muted-foreground">
                    {tipoLabel(m.tipo)}
                    {m.usuario.nombre_completo ? ` · ${m.usuario.nombre_completo}` : ''}
                  </p>
                </div>
                <div className="shrink-0 text-right">
                  <p className="text-sm font-medium">{m.cantidad}</p>
                  <p className="text-xs text-muted-foreground">{formatRelative(m.created_at)}</p>
                </div>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}
