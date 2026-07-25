'use client';

import { AlertTriangle, DollarSign, ShoppingCart, Wallet } from 'lucide-react';

import { KpiCard } from '@/components/dashboard/kpi-card';
import { TopFraganciasCard } from '@/components/dashboard/top-fragancias-card';
import { UltimosMovimientosCard } from '@/components/dashboard/ultimos-movimientos-card';
import { VentasPorMetodoCard } from '@/components/dashboard/ventas-por-metodo-card';
import { VentasPorVendedoraCard } from '@/components/dashboard/ventas-por-vendedora-card';
import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { formatCurrency } from '@/lib/formatters';
import { useAlertasStock } from '../api/use-alertas-stock';
import { useResumenHoy } from '../api/use-resumen-hoy';
import { useUltimosMovimientos } from '../api/use-ultimos-movimientos';

export function DashboardContent() {
  const resumen = useResumenHoy();
  const movimientos = useUltimosMovimientos();
  const alertas = useAlertasStock();

  if (resumen.isLoading) return <LoadingState />;
  if (resumen.error) return <ErrorState error={resumen.error} onRetry={() => resumen.refetch()} />;

  const data = resumen.data;
  if (!data) return null;

  const alertasTotal = alertas.data?.total ?? 0;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold">Dashboard</h1>
        <p className="text-sm text-muted-foreground">Resumen de operación del día</p>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <KpiCard icon={DollarSign} label="Ventas hoy" value={formatCurrency(data.total_dia)} />
        <KpiCard icon={ShoppingCart} label="Número de ventas" value={data.ventas_count.toString()} />
        <KpiCard
          icon={Wallet}
          label="Efectivo en caja"
          value={formatCurrency(data.por_metodo_pago.efectivo)}
        />
        <KpiCard
          icon={AlertTriangle}
          label="Alertas de stock"
          value={alertasTotal.toString()}
          variant={alertasTotal > 0 ? 'warning' : 'default'}
        />
      </div>

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <VentasPorMetodoCard totales={data.por_metodo_pago} totalDia={data.total_dia} />
        <UltimosMovimientosCard movimientos={movimientos.data} isLoading={movimientos.isLoading} />
      </div>

      <TopFraganciasCard fragancias={data.top_fragancias} />
      <VentasPorVendedoraCard vendedoras={data.por_vendedora} />
    </div>
  );
}
