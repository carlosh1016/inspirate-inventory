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
import { formatCurrency } from '@/lib/formatters';
import { useProducto } from '@/features/productos/api/use-producto';
import { useDeleteProducto } from '@/features/productos/api/use-delete-producto';
import { ProductoForm } from '@/features/productos/components/producto-form';
import { categoriaLabel } from '@/features/productos/types';

const BASE_PATH = '/inventario/productos';

export default function ProductoDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const { isAdmin } = usePermission();
  const { data: producto, isLoading, isError, error, refetch } = useProducto(Number(id));
  const deleteProducto = useDeleteProducto();
  const [editing, setEditing] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  if (isLoading) return <LoadingState />;
  if (isError || !producto) return <ErrorState error={error} onRetry={() => refetch()} />;

  if (editing) {
    return (
      <div className="max-w-2xl">
        <PageHeader title={`Editar: ${producto.nombre}`} backHref={BASE_PATH} backLabel="Volver a productos" />
        <ProductoForm initialData={producto} />
      </div>
    );
  }

  return (
    <div className="max-w-3xl">
      <PageHeader
        title={producto.nombre}
        description={`${categoriaLabel(producto.categoria)} · ${formatCurrency(producto.precio)}`}
        backHref={BASE_PATH}
        backLabel="Volver a productos"
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
              vitrina={producto.stock.vitrina}
              bodega={producto.stock.bodega}
              total={producto.stock.total}
              minimo={String(producto.stock_minimo)}
              unidad="unidades"
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Movimientos recientes</CardTitle>
          </CardHeader>
          <CardContent>
            <ItemMovimientosList tipoItem="producto" itemId={producto.id} unidad="unidades" />
          </CardContent>
        </Card>
      </div>

      <ConfirmDeleteDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title="Eliminar producto"
        description={`¿Seguro que deseas eliminar "${producto.nombre}"? Esta acción requiere que su stock esté en cero.`}
        onConfirm={async () => {
          await deleteProducto.mutateAsync(producto.id);
          toast.success('Producto eliminado');
          router.push(BASE_PATH);
        }}
      />
    </div>
  );
}
