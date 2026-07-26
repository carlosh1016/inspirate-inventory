'use client';

import { Suspense } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';

import { ErrorState } from '@/components/feedback/error-state';
import { PageHeader } from '@/components/page-header';
import { buttonVariants } from '@/components/ui/button';
import { usePermission } from '@/hooks/use-permission';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { useVentas } from '@/features/ventas/api/use-ventas';
import { VentasFiltersBar } from '@/features/ventas/components/ventas-filters';
import { VentasTable } from '@/features/ventas/components/ventas-table';
import type { VentasFilters } from '@/features/ventas/types';

const nuevaVentaButton = (
  <Link href="/nueva-venta" className={buttonVariants()}>
    <Plus className="size-4" />
    Nueva venta
  </Link>
);

function VentasView() {
  const router = useRouter();
  const { isAdmin } = usePermission();
  const { filters, setFilter } = useUrlFilters<VentasFilters>({
    defaults: {
      page: 1,
      metodo_pago_id: 0,
      usuario_id: 0,
      fecha_desde: '',
      fecha_hasta: '',
      con_descuento: false,
      total_min: '',
      total_max: '',
    },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      metodo_pago_id: (v) => Number(v) || 0,
      usuario_id: (v) => Number(v) || 0,
      fecha_desde: (v) => v ?? '',
      fecha_hasta: (v) => v ?? '',
      con_descuento: (v) => v === 'true',
      total_min: (v) => v ?? '',
      total_max: (v) => v ?? '',
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      metodo_pago_id: (v) => (v > 0 ? String(v) : null),
      usuario_id: (v) => (v > 0 ? String(v) : null),
      fecha_desde: (v) => v || null,
      fecha_hasta: (v) => v || null,
      con_descuento: (v) => (v ? 'true' : null),
      total_min: (v) => v || null,
      total_max: (v) => v || null,
    },
  });

  const { data, isLoading, isError, error, refetch } = useVentas(filters);

  return (
    <>
      <PageHeader title="Ventas" description="Historial de ventas registradas." action={nuevaVentaButton} />
      <VentasFiltersBar filters={filters} setFilter={setFilter} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <VentasTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          showVendedora={isAdmin}
          onPageChange={(p) => setFilter('page', p)}
          onRowClick={(v) => router.push(`/ventas/${v.id}`)}
        />
      )}
    </>
  );
}

export default function VentasPage() {
  return (
    <Suspense>
      <VentasView />
    </Suspense>
  );
}
