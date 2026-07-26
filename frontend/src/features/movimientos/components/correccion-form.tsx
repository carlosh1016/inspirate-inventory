'use client';

import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

import { useCorreccion } from '../api/use-correccion';
import { SingleMovimientoForm } from './single-movimiento-form';

export function CorreccionForm() {
  const router = useRouter();
  const correccion = useCorreccion();

  return (
    <SingleMovimientoForm
      cantidadLabel="Cantidad final que quedará en esta ubicación"
      cantidadHint="Es la cantidad absoluta que debe quedar, no la diferencia."
      allowZero
      submitLabel="Registrar corrección"
      submittingLabel="Registrando…"
      onSubmit={async (v) => {
        const result = await correccion.mutateAsync({
          tipo_item: v.tipo_item,
          item_id: v.item_id,
          ubicacion: v.ubicacion,
          cantidad_nueva: v.cantidad,
          motivo: v.motivo,
        });
        if (result.movimiento) {
          toast.success(
            `Corrección registrada. Stock: ${result.movimiento.stock_anterior} → ${result.movimiento.stock_posterior}`,
          );
        } else {
          toast.info(result.mensaje ?? 'El stock ya estaba en el valor solicitado.');
        }
        router.push('/inventario/movimientos');
      }}
    />
  );
}
