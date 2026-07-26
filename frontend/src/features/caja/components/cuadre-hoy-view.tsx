'use client';

import { useState } from 'react';
import { toast } from 'sonner';

import { ConfirmDeleteDialog } from '@/components/confirm-delete-dialog';
import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { Button } from '@/components/ui/button';
import { usePermission } from '@/hooks/use-permission';
import { formatCurrency, formatDate, formatRelative } from '@/lib/formatters';
import { cn } from '@/lib/utils';
import { useDeleteConsignacion, useDeletePago } from '../api/mutations';
import { useCuadreHoy } from '../api/use-cuadre-hoy';
import type { Consignacion, PagoCaja } from '../types';
import { CuadreAbrirDialog } from './cuadre-abrir-dialog';
import { CuadreAddConsignacionDialog } from './cuadre-add-consignacion-dialog';
import { CuadreAddPagoDialog } from './cuadre-add-pago-dialog';
import { CuadreCerrarDialog } from './cuadre-cerrar-dialog';
import { CuadreNoAbierto } from './cuadre-no-abierto';
import { CuadreSections } from './cuadre-sections';
import { CuadreWarningAnterior } from './cuadre-warning-anterior';

type DeleteTarget =
  | { kind: 'pago'; item: PagoCaja }
  | { kind: 'consignacion'; item: Consignacion }
  | null;

export function CuadreHoyView() {
  const { isAdmin } = usePermission();
  const { data: cuadre, isLoading, isError, error, refetch } = useCuadreHoy();

  const [abrirOpen, setAbrirOpen] = useState(false);
  const [cerrarOpen, setCerrarOpen] = useState(false);
  const [addPagoOpen, setAddPagoOpen] = useState(false);
  const [addConsignacionOpen, setAddConsignacionOpen] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<DeleteTarget>(null);

  const cuadreId = cuadre?.id ?? 0;
  const deletePago = useDeletePago(cuadreId);
  const deleteConsignacion = useDeleteConsignacion(cuadreId);

  if (isLoading) return <LoadingState />;
  if (isError) return <ErrorState error={error} onRetry={() => refetch()} />;

  if (!cuadre) {
    return (
      <>
        <CuadreWarningAnterior />
        <CuadreNoAbierto isAdmin={isAdmin} onAbrir={() => setAbrirOpen(true)} />
        <CuadreAbrirDialog open={abrirOpen} onOpenChange={setAbrirOpen} />
      </>
    );
  }

  const abierto = cuadre.estado === 'abierto';

  return (
    <>
      <CuadreWarningAnterior />

      <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-lg font-semibold">Caja · {formatDate(cuadre.fecha)}</h2>
          <p className="text-sm text-muted-foreground">Fondo inicial: {formatCurrency(cuadre.fondo_base)}</p>
        </div>
        <span
          className={cn(
            'rounded-md px-2.5 py-1 text-xs font-medium',
            abierto ? 'bg-success/10 text-success' : 'bg-muted text-muted-foreground',
          )}
        >
          {abierto ? 'ABIERTA' : 'CERRADA'}
        </span>
      </div>

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
        canDelete={isAdmin && abierto}
        onAddPago={() => setAddPagoOpen(true)}
        onAddConsignacion={() => setAddConsignacionOpen(true)}
        onDeletePago={(item) => setDeleteTarget({ kind: 'pago', item })}
        onDeleteConsignacion={(item) => setDeleteTarget({ kind: 'consignacion', item })}
      />

      {abierto && isAdmin && (
        <div className="mt-6 flex justify-end">
          <Button variant="destructive" onClick={() => setCerrarOpen(true)}>
            Cerrar caja del día
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
    </>
  );
}
