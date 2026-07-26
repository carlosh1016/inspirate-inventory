'use client';

import { use, useState } from 'react';
import { useRouter } from 'next/navigation';
import { MoreVertical, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

import { ConfirmDeleteDialog } from '@/components/confirm-delete-dialog';
import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { ItemMovimientosList } from '@/components/inventario/item-movimientos-list';
import { StockSummaryCard } from '@/components/inventario/stock-summary-card';
import { PageHeader } from '@/components/page-header';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { usePermission } from '@/hooks/use-permission';
import { useModelosLookup } from '@/features/modelos-envase/api/use-modelos-lookup';
import { useVarianteEnvase } from '@/features/variantes-envase/api/use-variante-envase';
import { useDeleteVarianteEnvase } from '@/features/variantes-envase/api/use-delete-variante-envase';
import { VarianteEnvaseForm } from '@/features/variantes-envase/components/variante-envase-form';

const BASE_PATH = '/inventario/envases/variantes';

export default function VarianteEnvaseDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const { isAdmin } = usePermission();
  const { data: variante, isLoading, isError, error, refetch } = useVarianteEnvase(Number(id));
  const { data: modelosMap } = useModelosLookup();
  const deleteVariante = useDeleteVarianteEnvase();
  const [editing, setEditing] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  if (isLoading) return <LoadingState />;
  if (isError || !variante) return <ErrorState error={error} onRetry={() => refetch()} />;

  const modeloLabel = modelosMap?.get(variante.modelo_envase_id) ?? `Modelo #${variante.modelo_envase_id}`;

  if (editing) {
    return (
      <div className="max-w-2xl">
        <PageHeader title={`Editar: ${variante.color}`} backHref={BASE_PATH} backLabel="Volver a variantes" />
        <VarianteEnvaseForm initialData={variante} />
      </div>
    );
  }

  return (
    <div className="max-w-3xl">
      <PageHeader
        title={`${modeloLabel} · ${variante.color}`}
        description={variante.activo ? undefined : 'Variante inactiva'}
        backHref={BASE_PATH}
        backLabel="Volver a variantes"
        action={
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={() => setEditing(true)}>
              Editar
            </Button>
            {isAdmin && (
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
            )}
          </div>
        }
      />

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Stock actual</CardTitle>
          </CardHeader>
          <CardContent>
            <StockSummaryCard
              vitrina={variante.stock.vitrina}
              bodega={variante.stock.bodega}
              total={variante.stock.total}
              minimo={String(variante.stock_minimo)}
              unidad="unidades"
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Movimientos recientes</CardTitle>
          </CardHeader>
          <CardContent>
            <ItemMovimientosList tipoItem="variante_envase" itemId={variante.id} unidad="unidades" />
          </CardContent>
        </Card>
      </div>

      <ConfirmDeleteDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="Eliminar variante"
        description={`¿Seguro que deseas eliminar la variante "${variante.color}"? Si tiene stock, ajústalo a cero primero.`}
        onConfirm={async () => {
          await deleteVariante.mutateAsync(variante.id);
          toast.success('Variante eliminada');
          router.push(BASE_PATH);
        }}
      />
    </div>
  );
}
