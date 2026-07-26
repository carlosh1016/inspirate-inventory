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
import { useProductos } from '@/features/productos/api/use-productos';
import { ProductosFiltersBar } from '@/features/productos/components/productos-filters';
import { ProductosTable } from '@/features/productos/components/productos-table';
import type { ProductosFilters } from '@/features/productos/types';

function ProductosView() {
  const router = useRouter();
  const { isAdmin } = usePermission();
  const { filters, setFilter } = useUrlFilters<ProductosFilters>({
    defaults: { page: 1, q: '', categoria: 'all', activo: 'true', stock_bajo: false },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      q: (v) => v ?? '',
      categoria: (v) => v ?? 'all',
      activo: (v) => v ?? 'true',
      stock_bajo: (v) => v === 'true',
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      q: (v) => v || null,
      categoria: (v) => (v === 'all' ? null : v),
      activo: (v) => (v === 'true' ? null : v),
      stock_bajo: (v) => (v ? 'true' : null),
    },
  });

  const { data, isLoading, isError, error, refetch } = useProductos(filters);

  const nuevoButton = isAdmin ? (
    <Link href="/inventario/productos/nuevo" className={buttonVariants()}>
      <Plus className="size-4" />
      Nuevo producto
    </Link>
  ) : null;

  return (
    <>
      <PageHeader
        title="Productos"
        description="Cremas, feromonas y otros productos del catálogo."
        action={nuevoButton}
      />
      <ProductosFiltersBar filters={filters} setFilter={setFilter} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <ProductosTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          onPageChange={(p) => setFilter('page', p)}
          onRowClick={(p) => router.push(`/inventario/productos/${p.id}`)}
          emptyAction={nuevoButton}
        />
      )}
    </>
  );
}

export default function ProductosPage() {
  return (
    <Suspense>
      <ProductosView />
    </Suspense>
  );
}
