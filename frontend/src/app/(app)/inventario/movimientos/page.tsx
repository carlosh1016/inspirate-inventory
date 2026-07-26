'use client';

import { Suspense } from 'react';

import { InventarioTabs } from '@/components/inventario/inventario-tabs';
import { RegistrarMovimientoMenu } from '@/components/inventario/registrar-movimiento-menu';
import { ErrorState } from '@/components/feedback/error-state';
import { PageHeader } from '@/components/page-header';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { useMovimientos, type MovimientosFilters } from '@/features/movimientos/api/use-movimientos';
import { MovimientosFiltersBar } from '@/features/movimientos/components/movimientos-filters';
import { MovimientosTable } from '@/features/movimientos/components/movimientos-table';

function MovimientosView() {
  const { filters, setFilter, setFilters } = useUrlFilters<MovimientosFilters>({
    defaults: {
      page: 1,
      tipo: 'all',
      tipo_item: 'all',
      ubicacion: 'all',
      item_id: 0,
      item_nombre: '',
      usuario_id: 0,
      fecha_desde: '',
      fecha_hasta: '',
    },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      tipo: (v) => v ?? 'all',
      tipo_item: (v) => v ?? 'all',
      ubicacion: (v) => v ?? 'all',
      item_id: (v) => Number(v) || 0,
      item_nombre: (v) => v ?? '',
      usuario_id: (v) => Number(v) || 0,
      fecha_desde: (v) => v ?? '',
      fecha_hasta: (v) => v ?? '',
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      tipo: (v) => (v === 'all' ? null : v),
      tipo_item: (v) => (v === 'all' ? null : v),
      ubicacion: (v) => (v === 'all' ? null : v),
      item_id: (v) => (v > 0 ? String(v) : null),
      item_nombre: (v) => v || null,
      usuario_id: (v) => (v > 0 ? String(v) : null),
      fecha_desde: (v) => v || null,
      fecha_hasta: (v) => v || null,
    },
  });

  const { data, isLoading, isError, error, refetch } = useMovimientos(filters);

  return (
    <>
      <PageHeader
        title="Movimientos"
        description="Historial de movimientos de inventario."
        action={<RegistrarMovimientoMenu />}
      />
      <InventarioTabs />
      <MovimientosFiltersBar filters={filters} setFilter={setFilter} setFilters={setFilters} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <MovimientosTable
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

export default function MovimientosPage() {
  return (
    <Suspense>
      <MovimientosView />
    </Suspense>
  );
}
