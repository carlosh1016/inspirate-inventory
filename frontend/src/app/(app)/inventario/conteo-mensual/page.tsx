'use client';

import { Suspense } from 'react';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';
import { toast } from 'sonner';

import { ErrorState } from '@/components/feedback/error-state';
import { RequireRole } from '@/components/guards/require-role';
import { InventarioTabs } from '@/components/inventario/inventario-tabs';
import { PageHeader } from '@/components/page-header';
import { Button } from '@/components/ui/button';
import { getErrorMessage } from '@/lib/errors';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { useCrearConteo } from '@/features/conteo-mensual/api/mutations';
import { useConteos } from '@/features/conteo-mensual/api/use-conteos';
import { ConteoTable } from '@/features/conteo-mensual/components/conteo-table';

function ConteoMensualView() {
  const router = useRouter();
  const { filters, setFilter } = useUrlFilters<{ page: number }>({
    defaults: { page: 1 },
    parsers: { page: (v) => Math.max(1, Number(v) || 1) },
    serializers: { page: (v) => (v === 1 ? null : String(v)) },
  });

  const { data, isLoading, isError, error, refetch } = useConteos(filters.page);
  const crear = useCrearConteo();

  const handleCrear = async () => {
    try {
      const conteo = await crear.mutateAsync();
      router.push(`/inventario/conteo-mensual/${conteo.id}`);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <>
      <PageHeader
        title="Conteo mensual"
        description="Compara el stock del sistema contra el conteo físico de cada fragancia, mes a mes."
        action={
          <Button onClick={handleCrear} disabled={crear.isPending}>
            <Plus className="size-4" />
            {crear.isPending ? 'Creando…' : 'Nuevo conteo de este mes'}
          </Button>
        }
      />
      <InventarioTabs />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <ConteoTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          onPageChange={(p) => setFilter('page', p)}
          onRowClick={(c) => router.push(`/inventario/conteo-mensual/${c.id}`)}
        />
      )}
    </>
  );
}

export default function ConteoMensualPage() {
  return (
    <RequireRole role="admin">
      <Suspense>
        <ConteoMensualView />
      </Suspense>
    </RequireRole>
  );
}
