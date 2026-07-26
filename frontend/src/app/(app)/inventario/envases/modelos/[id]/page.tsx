'use client';

import { use, useState } from 'react';
import { useRouter } from 'next/navigation';
import { MoreVertical, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

import { ConfirmDeleteDialog } from '@/components/confirm-delete-dialog';
import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { PageHeader } from '@/components/page-header';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { usePermission } from '@/hooks/use-permission';
import { formatCurrency } from '@/lib/formatters';
import { useModeloEnvase } from '@/features/modelos-envase/api/use-modelo-envase';
import { useDeleteModeloEnvase } from '@/features/modelos-envase/api/use-delete-modelo-envase';
import { ModeloEnvaseForm } from '@/features/modelos-envase/components/modelo-envase-form';

const BASE_PATH = '/inventario/envases/modelos';

export default function ModeloEnvaseDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const { isAdmin } = usePermission();
  const { data: modelo, isLoading, isError, error, refetch } = useModeloEnvase(Number(id));
  const deleteModelo = useDeleteModeloEnvase();
  const [editing, setEditing] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  if (isLoading) return <LoadingState />;
  if (isError || !modelo) return <ErrorState error={error} onRetry={() => refetch()} />;

  if (editing) {
    return (
      <div className="max-w-2xl">
        <PageHeader title={`Editar: ${modelo.tipo}`} backHref={BASE_PATH} backLabel="Volver a modelos" />
        <ModeloEnvaseForm initialData={modelo} />
      </div>
    );
  }

  const rows: { label: string; value: string }[] = [
    { label: 'Tamaño', value: `${modelo.tamano_oz} oz` },
    { label: 'Equivalencia', value: `${modelo.equiv_gramos} g` },
    { label: 'Precio solo', value: formatCurrency(modelo.precio_solo) },
    { label: 'Precio con fragancia', value: formatCurrency(modelo.precio_con_fragancia) },
    { label: 'Precio recarga', value: formatCurrency(modelo.precio_recarga) },
    { label: 'Variantes activas', value: String(modelo.variantes_activas) },
  ];

  return (
    <div className="max-w-2xl">
      <PageHeader
        title={modelo.tipo}
        description={modelo.activo ? undefined : 'Modelo inactivo'}
        backHref={BASE_PATH}
        backLabel="Volver a modelos"
        action={
          isAdmin && (
            <div className="flex items-center gap-2">
              <Button variant="outline" onClick={() => setEditing(true)}>
                Editar
              </Button>
              <DropdownMenu>
                <DropdownMenuTrigger
                  aria-label="Más acciones"
                  className="inline-flex size-8 items-center justify-center rounded-md border border-border hover:bg-muted"
                >
                  <MoreVertical className="size-4" />
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => setConfirmOpen(true)}>
                    <Trash2 className="size-4" />
                    Eliminar
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          )
        }
      />

      <Card>
        <CardContent>
          <dl className="divide-y divide-border">
            {rows.map((row) => (
              <div key={row.label} className="flex items-center justify-between py-2.5 text-sm">
                <dt className="text-muted-foreground">{row.label}</dt>
                <dd className="font-medium tabular-nums">{row.value}</dd>
              </div>
            ))}
          </dl>
        </CardContent>
      </Card>

      <ConfirmDeleteDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="Eliminar modelo de envase"
        description={
          modelo.variantes_activas > 0
            ? `Este modelo tiene ${modelo.variantes_activas} variante(s) activa(s). Elimínalas o desactívalas primero.`
            : `¿Seguro que deseas eliminar el modelo "${modelo.tipo}"?`
        }
        onConfirm={async () => {
          await deleteModelo.mutateAsync(modelo.id);
          toast.success('Modelo eliminado');
          router.push(BASE_PATH);
        }}
      />
    </div>
  );
}
