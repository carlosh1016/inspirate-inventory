'use client';

import { Suspense } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';

import { PageHeader } from '@/components/page-header';
import { buttonVariants } from '@/components/ui/button';
import { ErrorState } from '@/components/feedback/error-state';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { FraganciasFiltersBar } from '@/features/fragancias/components/fragancias-filters';
import { FraganciasTable } from '@/features/fragancias/components/fragancias-table';
import { useFragancias } from '@/features/fragancias/api/use-fragancias';
import type { FraganciasFilters } from '@/features/fragancias/types';

const NuevaButton = (
  <Link href="/inventario/fragancias/nueva" className={buttonVariants()}>
    <Plus className="size-4" />
    Nueva fragancia
  </Link>
);

function FraganciasView() {
  const router = useRouter();
  const { filters, setFilter } = useUrlFilters<FraganciasFilters>({
    defaults: { page: 1, q: '', genero: 'all', activo: 'true', stock_bajo: false },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      q: (v) => v ?? '',
      genero: (v) => v ?? 'all',
      activo: (v) => v ?? 'true',
      stock_bajo: (v) => v === 'true',
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      q: (v) => v || null,
      genero: (v) => (v === 'all' ? null : v),
      activo: (v) => (v === 'true' ? null : v),
      stock_bajo: (v) => (v ? 'true' : null),
    },
  });

  const { data, isLoading, isError, error, refetch } = useFragancias(filters);

  return (
    <>
      <PageHeader
        title="Fragancias"
        description="Catálogo de fragancias con stock por sede."
        action={NuevaButton}
      />
      <FraganciasFiltersBar filters={filters} setFilter={setFilter} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <FraganciasTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          onPageChange={(p) => setFilter('page', p)}
          onRowClick={(f) => router.push(`/inventario/fragancias/${f.id}`)}
          emptyAction={NuevaButton}
        />
      )}
    </>
  );
}

export default function FraganciasPage() {
  return (
    <Suspense>
      <FraganciasView />
    </Suspense>
  );
}
