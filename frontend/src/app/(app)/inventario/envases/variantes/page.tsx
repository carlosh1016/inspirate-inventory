'use client';

import { Suspense } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';

import { EnvasesTabs } from '@/components/inventario/envases-tabs';
import { ErrorState } from '@/components/feedback/error-state';
import { PageHeader } from '@/components/page-header';
import { buttonVariants } from '@/components/ui/button';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { useModelosLookup } from '@/features/modelos-envase/api/use-modelos-lookup';
import { useVariantesEnvase } from '@/features/variantes-envase/api/use-variantes-envase';
import { VariantesEnvaseFiltersBar } from '@/features/variantes-envase/components/variantes-envase-filters';
import { VariantesEnvaseTable } from '@/features/variantes-envase/components/variantes-envase-table';
import type { VariantesEnvaseFilters } from '@/features/variantes-envase/types';

const nuevaButton = (
  <Link href="/inventario/envases/variantes/nueva" className={buttonVariants()}>
    <Plus className="size-4" />
    Nueva variante
  </Link>
);

function VariantesView() {
  const router = useRouter();
  const { data: modelosMap } = useModelosLookup();
  const { filters, setFilter } = useUrlFilters<VariantesEnvaseFilters>({
    defaults: { page: 1, q: '', modelo_envase_id: 0, activo: 'true', stock_bajo: false },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      q: (v) => v ?? '',
      modelo_envase_id: (v) => Number(v) || 0,
      activo: (v) => v ?? 'true',
      stock_bajo: (v) => v === 'true',
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      q: (v) => v || null,
      modelo_envase_id: (v) => (v > 0 ? String(v) : null),
      activo: (v) => (v === 'true' ? null : v),
      stock_bajo: (v) => (v ? 'true' : null),
    },
  });

  const { data, isLoading, isError, error, refetch } = useVariantesEnvase(filters);

  return (
    <>
      <PageHeader
        title="Variantes de envase"
        description="Envases por color con su stock."
        action={nuevaButton}
      />
      <EnvasesTabs />
      <VariantesEnvaseFiltersBar filters={filters} setFilter={setFilter} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <VariantesEnvaseTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          modelosMap={modelosMap}
          onPageChange={(p) => setFilter('page', p)}
          onRowClick={(v) => router.push(`/inventario/envases/variantes/${v.id}`)}
          emptyAction={nuevaButton}
        />
      )}
    </>
  );
}

export default function VariantesEnvasePage() {
  return (
    <Suspense>
      <VariantesView />
    </Suspense>
  );
}
