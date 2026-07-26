'use client';

import { use, useState } from 'react';
import { toast } from 'sonner';

import { ConfirmDeleteDialog } from '@/components/confirm-delete-dialog';
import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { Button } from '@/components/ui/button';
import { formatCurrency, formatDate, formatRelative } from '@/lib/formatters';
import { cn } from '@/lib/utils';
import { useDeleteConsignacion, useDeletePago } from '@/features/caja/api/mutations';
import { useCuadre } from '@/features/caja/api/use-cuadre';
import { CuadreAddConsignacionDialog } from '@/features/caja/components/cuadre-add-consignacion-dialog';
import { CuadreAddPagoDialog } from '@/features/caja/components/cuadre-add-pago-dialog';
import { CuadreCerrarDialog } from '@/features/caja/components/cuadre-cerrar-dialog';
import { CuadreSections } from '@/features/caja/components/cuadre-sections';
import type { Consignacion, PagoCaja } from '@/features/caja/types';

type DeleteTarget =
  | { kind: 'pago'; item: PagoCaja }
  | { kind: 'consignacion'; item: Consignacion }
  | null;

function CuadreDetalle({ cuadreId }: { cuadreId: number }) {
  const { data: cuadre, isLoading, isError, error, refetch } = useCuadre(cuadreId);
  const deletePago = useDeletePago(cuadreId);
  const deleteConsignacion = useDeleteConsignacion(cuadreId);
  const [cerrarOpen, setCerrarOpen] = useState(false);
  const [addPagoOpen, setAddPagoOpen] = useState(false);
  const [addConsignacionOpen, setAddConsignacionOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<DeleteTarget>(null);

  if (isLoading) return <LoadingState />;
  if (isError || !cuadre) return <ErrorState error={error} onRetry={() => refetch()} />;

  const abierto = cuadre.estado === 'abierto';

  return (
    <div className="max-w-3xl">
      <PageHeader
        title={`Caja · ${formatDate(cuadre.fecha)}`}
        description={`Fondo inicial: ${formatCurrency(cuadre.fondo_base)}`}
        backHref="/caja/historial"
        backLabel="Volver al historial"
        action={
          <span
            className={cn(
              'rounded-md px-2.5 py-1 text-xs font-medium',
              abierto ? 'bg-success/10 text-success' : 'bg-muted text-muted-foreground',
            )}
          >
            {abierto ? 'ABIERTA' : 'CERRADA'}
          </span>
        }
      />

      {!abierto && cuadre.cerrado_por && (
        <p className="mb-4 text-sm text-muted-foreground">
          Cerrado por {cuadre.cerrado_por.nombre_completo}
          {cuadre.cerrado_at ? ` ${formatRelative(cuadre.cerrado_at)}` : ''}.
          {Number.parseFloat(cuadre.valor_turno) > 0 && ` Valor turno: ${formatCurrency(cuadre.valor_turno)}.`}
          {cuadre.observaciones ? ` ${cuadre.observaciones}` : ''}
        </p>
      )}

      <CuadreSections
        cuadre={cuadre}
        canManage={abierto}
        canDelete={abierto}
        onAddPago={() => setAddPagoOpen(true)}
        onAddConsignacion={() => setAddConsignacionOpen(true)}
        onDeletePago={(item) => setDeleteTarget({ kind: 'pago', item })}
        onDeleteConsignacion={(item) => setDeleteTarget({ kind: 'consignacion', item })}
      />

      {abierto && (
        <div className="mt-6 flex justify-end">
          <Button variant="destructive" onClick={() => setCerrarOpen(true)}>
            Cerrar caja
          </Button>
        </div>
      )}

      <CuadreAddPagoDialog cuadreId={cuadre.id} open={addPagoOpen} onOpenChange={setAddPagoOpen} />
      <CuadreAddConsignacionDialog cuadreId={cuadre.id} open={addConsignacionOpen} onOpenChange={setAddConsignacionOpen} />
      <CuadreCerrarDialog cuadre={cuadre} open={cerrarOpen} onOpenChange={setCerrarOpen} />

      <ConfirmDeleteDialog
        open={deleteTarget !== null}
        onOpenChange={(open) => !open && setDeleteTarget(null)}
        title={deleteTarget?.kind === 'pago' ? 'Eliminar pago' : 'Eliminar consignación'}
        description="¿Seguro que deseas eliminar este registro de la caja?"
        onConfirm={async () => {
          if (!deleteTarget) return;
          if (deleteTarget.kind === 'pago') {
            await deletePago.mutateAsync(deleteTarget.item.id);
          } else {
            await deleteConsignacion.mutateAsync(deleteTarget.item.id);
          }
          toast.success('Registro eliminado');
          setDeleteTarget(null);
        }}
      />
    </div>
  );
}

export default function CuadreDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return (
    <RequireRole role="admin">
      <CuadreDetalle cuadreId={Number(id)} />
    </RequireRole>
  );
}
