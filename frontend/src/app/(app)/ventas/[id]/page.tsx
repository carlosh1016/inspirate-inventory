'use client';

import { use } from 'react';
import Link from 'next/link';

import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { PageHeader } from '@/components/page-header';
import { ConsecutivoBadge } from '@/components/ventas/consecutivo-badge';
import { MetodoPagoBadge } from '@/components/ventas/metodo-pago-badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { formatCurrency, formatDateTime } from '@/lib/formatters';
import { TIPO_ITEM_LABEL } from '@/features/movimientos/tipo-meta';
import { useVenta } from '@/features/ventas/api/use-venta';
import { useVentaMovimientos } from '@/features/ventas/api/use-venta-movimientos';
import { VentaItemsList } from '@/features/ventas/components/venta-items-list';
import { VentaObservacionesInline } from '@/features/ventas/components/venta-observaciones-inline';

export default function VentaDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const ventaId = Number(id);
  const { data: venta, isLoading, isError, error, refetch } = useVenta(ventaId);
  const { data: movimientos } = useVentaMovimientos(ventaId);

  if (isLoading) return <LoadingState />;
  if (isError || !venta) return <ErrorState error={error} onRetry={() => refetch()} />;

  const conDescuento = Number.parseFloat(venta.descuento_monto) > 0;

  return (
    <div className="max-w-3xl">
      <PageHeader
        title={`Venta ${venta.consecutivo}`}
        description={`Registrada el ${formatDateTime(venta.created_at)} · Por ${venta.usuario.nombre_completo}`}
        backHref="/ventas"
        backLabel="Volver a ventas"
        action={<ConsecutivoBadge value={venta.consecutivo} />}
      />

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Ítems</CardTitle>
          </CardHeader>
          <CardContent>
            <VentaItemsList items={venta.items} />
          </CardContent>
        </Card>

        <Card>
          <CardContent className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Subtotal</span>
              <span className="tabular-nums">{formatCurrency(venta.subtotal)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Descuento</span>
              <span className={conDescuento ? 'text-success tabular-nums' : 'tabular-nums'}>
                {conDescuento ? `-${formatCurrency(venta.descuento_monto)} (${venta.descuento_pct}%)` : formatCurrency(0)}
              </span>
            </div>
            <div className="flex items-baseline justify-between border-t border-border pt-2">
              <span className="font-medium">Total</span>
              <span className="text-xl font-semibold tabular-nums">{formatCurrency(venta.total)}</span>
            </div>
            <div className="flex justify-between pt-2">
              <span className="text-muted-foreground">Método de pago</span>
              <MetodoPagoBadge nombre={venta.metodo_pago.nombre} codigo={venta.metodo_pago.codigo} />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Observaciones</CardTitle>
          </CardHeader>
          <CardContent>
            <VentaObservacionesInline ventaId={venta.id} observaciones={venta.observaciones} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Movimientos generados</CardTitle>
          </CardHeader>
          <CardContent>
            {!movimientos || movimientos.length === 0 ? (
              <p className="text-sm text-muted-foreground">Sin movimientos asociados.</p>
            ) : (
              <ul className="space-y-1.5 text-sm">
                {movimientos.map((mov) => (
                  <li key={mov.id}>
                    <Link
                      href={`/inventario/movimientos/${mov.id}`}
                      className="flex items-center justify-between gap-3 rounded-md px-1 py-1 hover:bg-muted/50"
                    >
                      <span>
                        {TIPO_ITEM_LABEL[mov.tipo_item]} {mov.item.nombre} · {mov.ubicacion}
                      </span>
                      <span
                        className={
                          Number.parseFloat(mov.cantidad) < 0 ? 'text-destructive tabular-nums' : 'text-success tabular-nums'
                        }
                      >
                        {mov.cantidad}
                      </span>
                    </Link>
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
