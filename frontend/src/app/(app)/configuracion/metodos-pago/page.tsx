'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { ErrorState } from '@/components/feedback/error-state';
import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { buttonVariants } from '@/components/ui/button';
import { useMetodosPago } from '@/features/metodos-pago/api/use-metodos-pago';
import type { MetodoPago } from '@/features/metodos-pago/types';

const nuevoButton = (
  <Link href="/configuracion/metodos-pago/nuevo" className={buttonVariants()}>
    <Plus className="size-4" />
    Nuevo método
  </Link>
);

const columns: Column<MetodoPago>[] = [
  {
    key: 'nombre',
    header: 'Nombre',
    cell: (m) => (
      <span className="font-medium">
        {m.nombre}
        {!m.activo && <span className="ml-2 text-xs text-muted-foreground">(inactivo)</span>}
      </span>
    ),
  },
  { key: 'codigo', header: 'Código', cell: (m) => <span className="text-muted-foreground">{m.codigo}</span> },
];

export default function MetodosPagoPage() {
  const router = useRouter();
  const { data, isLoading, isError, error, refetch } = useMetodosPago();

  return (
    <RequireRole role="admin">
      <PageHeader
        title="Métodos de pago"
        description="Formas de pago disponibles en las ventas."
        action={nuevoButton}
      />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <DataTable
          columns={columns}
          data={data ?? []}
          isLoading={isLoading}
          emptyMessage="Aún no hay métodos de pago."
          emptyAction={nuevoButton}
          rowKey={(m) => m.id}
          onRowClick={(m) => router.push(`/configuracion/metodos-pago/${m.id}`)}
        />
      )}
    </RequireRole>
  );
}
