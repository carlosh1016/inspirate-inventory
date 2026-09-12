'use client';

import { use, useState } from 'react';
import { Download } from 'lucide-react';
import { toast } from 'sonner';

import { ConfirmDeleteDialog } from '@/components/confirm-delete-dialog';
import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { Button } from '@/components/ui/button';
import { getErrorMessage } from '@/lib/errors';
import { formatRelative } from '@/lib/formatters';
import { cn } from '@/lib/utils';
import { useCerrarConteo, useExportarConteoInventario } from '@/features/conteo-mensual/api/mutations';
import { useConteo } from '@/features/conteo-mensual/api/use-conteo';
import { ConteoItemsTable } from '@/features/conteo-mensual/components/conteo-items-table';
import { formatPeriodo } from '@/features/conteo-mensual/format-periodo';

function ConteoDetalle({ conteoId }: { conteoId: number }) {
  const { data: conteo, isLoading, isError, error, refetch } = useConteo(conteoId);
  const cerrar = useCerrarConteo(conteoId);
  const exportar = useExportarConteoInventario();
  const [cerrarOpen, setCerrarOpen] = useState(false);

  if (isLoading) return <LoadingState />;
  if (isError || !conteo) return <ErrorState error={error} onRetry={() => refetch()} />;

  const abierto = conteo.estado === 'abierto';

  return (
    <div className="max-w-4xl">
      <PageHeader
        title={`Conteo mensual · ${formatPeriodo(conteo.periodo)}`}
        description={`Creado por ${conteo.creado_por.nombre_completo}`}
        backHref="/inventario/conteo-mensual"
        backLabel="Volver a conteos"
        action={
          <span
            className={cn(
              'rounded-md px-2.5 py-1 text-xs font-medium',
              abierto ? 'bg-success/10 text-success' : 'bg-muted text-muted-foreground',
            )}
          >
            {abierto ? 'ABIERTO' : 'CERRADO'}
          </span>
        }
      />

      {!abierto && conteo.cerrado_por && (
        <p className="mb-4 text-sm text-muted-foreground">
          Cerrado por {conteo.cerrado_por.nombre_completo}
          {conteo.cerrado_at ? ` ${formatRelative(conteo.cerrado_at)}` : ''}.
        </p>
      )}

      <ConteoItemsTable conteoId={conteo.id} items={conteo.items} disabled={!abierto} />

      <div className="mt-6 flex justify-end gap-2">
        <Button
          variant="outline"
          disabled={exportar.isPending}
          onClick={async () => {
            try {
              await exportar.mutateAsync({ id: conteo.id, periodo: conteo.periodo });
            } catch (err) {
              toast.error(getErrorMessage(err));
            }
          }}
        >
          <Download className="size-4" />
          {exportar.isPending ? 'Exportando…' : 'Exportar Excel'}
        </Button>
        {abierto && (
          <Button variant="destructive" onClick={() => setCerrarOpen(true)}>
            Cerrar conteo
          </Button>
        )}
      </div>

      <ConfirmDeleteDialog
        open={cerrarOpen}
        onOpenChange={setCerrarOpen}
        title="Cerrar conteo mensual"
        description="Una vez cerrado, no podrás editar los gramos físicos de este conteo. El stock del sistema no se modifica al cerrar — si encontraste una diferencia, corrígela por separado en Inventario › Movimientos › Corrección."
        confirmLabel="Cerrar conteo definitivamente"
        onConfirm={async () => {
          await cerrar.mutateAsync();
          toast.success('Conteo cerrado correctamente');
        }}
      />
    </div>
  );
}

export default function ConteoMensualDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return (
    <RequireRole role="admin">
      <ConteoDetalle conteoId={Number(id)} />
    </RequireRole>
  );
}
