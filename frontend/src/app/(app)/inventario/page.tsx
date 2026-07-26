'use client';

import { Suspense } from 'react';

import { InventarioTabs } from '@/components/inventario/inventario-tabs';
import { RegistrarMovimientoMenu } from '@/components/inventario/registrar-movimiento-menu';
import { ErrorState } from '@/components/feedback/error-state';
import { PageHeader } from '@/components/page-header';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { useStock } from '@/features/stock/api/use-stock';
import { StockFiltersBar } from '@/features/stock/components/stock-filters';
import { StockTable } from '@/features/stock/components/stock-table';
import type { StockFilters } from '@/features/stock/types';

function InventarioView() {
  const { filters, setFilter } = useUrlFilters<StockFilters>({
    defaults: { page: 1, tipo_item: 'all', ubicacion: 'all', stock_bajo: false, stock_cero: false },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      tipo_item: (v) => v ?? 'all',
      ubicacion: (v) => v ?? 'all',
      stock_bajo: (v) => v === 'true',
      stock_cero: (v) => v === 'true',
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      tipo_item: (v) => (v === 'all' ? null : v),
      ubicacion: (v) => (v === 'all' ? null : v),
      stock_bajo: (v) => (v ? 'true' : null),
      stock_cero: (v) => (v ? 'true' : null),
    },
  });

  const { data, isLoading, isError, error, refetch } = useStock(filters);

  return (
    <>
      <PageHeader
        title="Inventario"
        description="Vista general del stock actual."
        action={<RegistrarMovimientoMenu />}
      />
      <InventarioTabs />
      <StockFiltersBar filters={filters} setFilter={setFilter} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <StockTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          onPageChange={(p) => setFilter('page', p)}
        />
      )}
    </>
  );
}

export default function InventarioPage() {
  return (
    <Suspense>
      <InventarioView />
    </Suspense>
  );
}
