'use client';

import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

import { useDanado } from '../api/use-danado';
import { SingleMovimientoForm } from './single-movimiento-form';

export function DanadoForm() {
  const router = useRouter();
  const danado = useDanado();

  return (
    <SingleMovimientoForm
      cantidadLabel="Cantidad dañada"
      cantidadHint="Se descontará del stock de la ubicación seleccionada."
      submitLabel="Registrar dañado"
      submittingLabel="Registrando…"
      onSubmit={async (v) => {
        await danado.mutateAsync({
          tipo_item: v.tipo_item,
          item_id: v.item_id,
          ubicacion: v.ubicacion,
          cantidad: v.cantidad,
          motivo: v.motivo,
        });
        toast.success('Producto dañado registrado');
        router.push('/inventario/movimientos');
      }}
    />
  );
}
