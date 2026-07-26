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
import { useFragancia } from '@/features/fragancias/api/use-fragancia';
import { useDeleteFragancia } from '@/features/fragancias/api/use-delete-fragancia';
import { FraganciaForm } from '@/features/fragancias/components/fragancia-form';

export default function FraganciaDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const fraganciaId = Number(id);
  const router = useRouter();
  const { isAdmin } = usePermission();
  const { data: fragancia, isLoading, isError, error, refetch } = useFragancia(fraganciaId);
  const deleteFragancia = useDeleteFragancia();
  const [editing, setEditing] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  if (isLoading) return <LoadingState />;
  if (isError || !fragancia) return <ErrorState error={error} onRetry={() => refetch()} />;

  if (editing) {
    return (
      <div className="max-w-2xl">
        <PageHeader
          title={`Editar: ${fragancia.nombre_comercial}`}
          backHref="/inventario/fragancias"
          backLabel="Volver a fragancias"
        />
        <FraganciaForm initialData={fragancia} />
      </div>
    );
  }

  const subtitle = [fragancia.nombre_alternativo, fragancia.genero]
    .filter(Boolean)
    .join(' · ');

  return (
    <div className="max-w-3xl">
      <PageHeader
        title={fragancia.nombre_comercial}
        description={subtitle}
        backHref="/inventario/fragancias"
        backLabel="Volver a fragancias"
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
              vitrina={fragancia.stock.vitrina}
              bodega={fragancia.stock.bodega}
              total={fragancia.stock.total}
              minimo={fragancia.gramos_minimo}
              unidad="gramos"
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Movimientos recientes</CardTitle>
          </CardHeader>
          <CardContent>
            <ItemMovimientosList tipoItem="fragancia" itemId={fragancia.id} unidad="gramos" />
          </CardContent>
        </Card>
      </div>

      <ConfirmDeleteDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="Eliminar fragancia"
        description={`¿Seguro que deseas eliminar "${fragancia.nombre_comercial}"? Esta acción requiere que su stock esté en cero.`}
        onConfirm={async () => {
          await deleteFragancia.mutateAsync(fragancia.id);
          toast.success('Fragancia eliminada');
          router.push('/inventario/fragancias');
        }}
      />
    </div>
  );
}
