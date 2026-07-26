'use client';

import { use, useState } from 'react';
import { useRouter } from 'next/navigation';
import { MoreVertical, Trash2 } from 'lucide-react';
import { toast } from 'sonner';

import { ConfirmDeleteDialog } from '@/components/confirm-delete-dialog';
import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useMetodoPago } from '@/features/metodos-pago/api/use-metodo-pago';
import { useDeleteMetodoPago } from '@/features/metodos-pago/api/use-delete-metodo-pago';
import { MetodoPagoForm } from '@/features/metodos-pago/components/metodo-pago-form';

const BASE_PATH = '/configuracion/metodos-pago';

export default function MetodoPagoDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const { data: metodo, isLoading, isError, error, refetch } = useMetodoPago(Number(id));
  const deleteMetodo = useDeleteMetodoPago();
  const [confirmOpen, setConfirmOpen] = useState(false);

  return (
    <RequireRole role="admin">
      {isLoading ? (
        <LoadingState />
      ) : isError || !metodo ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <div className="max-w-lg">
          <PageHeader
            title={`Editar: ${metodo.nombre}`}
            backHref={BASE_PATH}
            backLabel="Volver a métodos de pago"
            action={
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
            }
          />
          <MetodoPagoForm initialData={metodo} />

          <ConfirmDeleteDialog
            open={confirmOpen}
            onOpenChange={setConfirmOpen}
            title="Eliminar método de pago"
            description={`¿Seguro que deseas eliminar "${metodo.nombre}"?`}
            onConfirm={async () => {
              await deleteMetodo.mutateAsync(metodo.id);
              toast.success('Método de pago eliminado');
              router.push(BASE_PATH);
            }}
          />
        </div>
      )}
    </RequireRole>
  );
}
