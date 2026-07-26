'use client';

import { Suspense } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';

import { EnvasesTabs } from '@/components/inventario/envases-tabs';
import { PageHeader } from '@/components/page-header';
import { ErrorState } from '@/components/feedback/error-state';
import { buttonVariants } from '@/components/ui/button';
import { usePermission } from '@/hooks/use-permission';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { useModelosEnvase } from '@/features/modelos-envase/api/use-modelos-envase';
import { ModelosEnvaseFiltersBar } from '@/features/modelos-envase/components/modelos-envase-filters';
import { ModelosEnvaseTable } from '@/features/modelos-envase/components/modelos-envase-table';
import type { ModelosEnvaseFilters } from '@/features/modelos-envase/types';

function ModelosView() {
  const router = useRouter();
  const { isAdmin } = usePermission();
  const { filters, setFilter } = useUrlFilters<ModelosEnvaseFilters>({
    defaults: { page: 1, q: '', activo: 'true' },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      q: (v) => v ?? '',
      activo: (v) => v ?? 'true',
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      q: (v) => v || null,
      activo: (v) => (v === 'true' ? null : v),
    },
  });

  const { data, isLoading, isError, error, refetch } = useModelosEnvase(filters);

  const nuevoButton = isAdmin ? (
    <Link href="/inventario/envases/modelos/nuevo" className={buttonVariants()}>
      <Plus className="size-4" />
      Nuevo modelo
    </Link>
  ) : null;

  return (
    <>
      <PageHeader
        title="Modelos de envase"
        description="Tipos de envase con sus tamaños y precios."
        action={nuevoButton}
      />
      <EnvasesTabs />
      <ModelosEnvaseFiltersBar filters={filters} setFilter={setFilter} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <ModelosEnvaseTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          onPageChange={(p) => setFilter('page', p)}
          onRowClick={(m) => router.push(`/inventario/envases/modelos/${m.id}`)}
          emptyAction={nuevoButton}
        />
      )}
    </>
  );
}

export default function ModelosEnvasePage() {
  return (
    <Suspense>
      <ModelosView />
    </Suspense>
  );
}
