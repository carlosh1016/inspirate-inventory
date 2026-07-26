'use client';

import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

import { useAjuste } from '../api/use-ajuste';
import { SingleMovimientoForm } from './single-movimiento-form';

export function AjusteForm() {
  const router = useRouter();
  const ajuste = useAjuste();

  return (
    <SingleMovimientoForm
      cantidadLabel="Cantidad final que quedará en esta ubicación"
      cantidadHint="Es la cantidad absoluta que debe quedar, no la diferencia."
      allowZero
      submitLabel="Registrar ajuste"
      submittingLabel="Registrando…"
      onSubmit={async (v) => {
        const result = await ajuste.mutateAsync({
          tipo_item: v.tipo_item,
          item_id: v.item_id,
          ubicacion: v.ubicacion,
          cantidad_nueva: v.cantidad,
          motivo: v.motivo,
        });
        if (result.movimiento) {
          toast.success(
            `Ajuste registrado. Stock: ${result.movimiento.stock_anterior} → ${result.movimiento.stock_posterior}`,
          );
        } else {
          toast.info(result.mensaje ?? 'El stock ya estaba en el valor solicitado.');
        }
        router.push('/inventario/movimientos');
      }}
    />
  );
}
