'use client';

import { Suspense } from 'react';
import { useRouter } from 'next/navigation';

import { ErrorState } from '@/components/feedback/error-state';
import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { useCuadres } from '@/features/caja/api/use-cuadres';
import { CuadresFiltersBar } from '@/features/caja/components/cuadres-filters';
import { CuadresTable } from '@/features/caja/components/cuadres-table';
import type { CuadresFilters } from '@/features/caja/types';

function HistorialView() {
  const router = useRouter();
  const { filters, setFilter } = useUrlFilters<CuadresFilters>({
    defaults: { page: 1, estado: 'all', fecha_desde: '', fecha_hasta: '' },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      estado: (v) => v ?? 'all',
      fecha_desde: (v) => v ?? '',
      fecha_hasta: (v) => v ?? '',
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      estado: (v) => (v === 'all' ? null : v),
      fecha_desde: (v) => v || null,
      fecha_hasta: (v) => v || null,
    },
  });

  const { data, isLoading, isError, error, refetch } = useCuadres(filters);

  return (
    <>
      <PageHeader
        title="Historial de cuadres"
        description="Cuadres de caja de días anteriores."
        backHref="/caja"
        backLabel="Volver a caja"
      />
      <CuadresFiltersBar filters={filters} setFilter={setFilter} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <CuadresTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          onPageChange={(p) => setFilter('page', p)}
          onRowClick={(c) => router.push(`/caja/${c.id}`)}
        />
      )}
    </>
  );
}

export default function CajaHistorialPage() {
  return (
    <RequireRole role="admin">
      <Suspense>
        <HistorialView />
      </Suspense>
    </RequireRole>
  );
}
