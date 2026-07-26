'use client';

import { Suspense } from 'react';
import Link from 'next/link';
import { ArrowRightLeft } from 'lucide-react';

import { InventarioTabs } from '@/components/inventario/inventario-tabs';
import { ErrorState } from '@/components/feedback/error-state';
import { PageHeader } from '@/components/page-header';
import { buttonVariants } from '@/components/ui/button';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { useAlertasStock } from '@/features/stock/api/use-alertas-stock';
import { StockTable } from '@/features/stock/components/stock-table';

type AlertasFilters = { page: number };

function AlertasView() {
  const { filters, setFilter } = useUrlFilters<AlertasFilters>({
    defaults: { page: 1 },
    parsers: { page: (v) => Math.max(1, Number(v) || 1) },
    serializers: { page: (v) => (v === 1 ? null : String(v)) },
  });

  const { data, isLoading, isError, error, refetch } = useAlertasStock(filters.page);

  return (
    <>
      <PageHeader title="Alertas de stock" description="Ítems por debajo de su mínimo." />
      <InventarioTabs />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <StockTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          onPageChange={(p) => setFilter('page', p)}
          emptyMessage="Todo el inventario está en niveles saludables."
          renderActions={(row) =>
            Number.parseFloat(row.stock_bodega) > 0 ? (
              <Link
                href={`/inventario/movimientos/traslado?tipo_item=${row.tipo_item}&item_id=${row.item_id}&nombre=${encodeURIComponent(row.nombre)}`}
                className={buttonVariants({ variant: 'outline', size: 'sm' })}
              >
                <ArrowRightLeft className="size-3.5" />
                Traslado desde bodega
              </Link>
            ) : null
          }
        />
      )}
    </>
  );
}

export default function AlertasPage() {
  return (
    <Suspense>
      <AlertasView />
    </Suspense>
  );
}
