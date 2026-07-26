'use client';

import { Suspense } from 'react';
import { useSearchParams } from 'next/navigation';

import { PageHeader } from '@/components/page-header';
import { TrasladoForm } from '@/features/movimientos/components/traslado-form';
import {
  emptyMovimientoItem,
  type MovimientoItemInput,
} from '@/features/movimientos/components/movimiento-item-form';
import type { TipoItem } from '@/types/domain';

const VALID_TIPOS: TipoItem[] = ['fragancia', 'variante_envase', 'producto'];

function TrasladoView() {
  const params = useSearchParams();
  const tipoItem = params.get('tipo_item');
  const itemId = Number(params.get('item_id')) || 0;
  const nombre = params.get('nombre') ?? '';

  let initialItem: MovimientoItemInput | undefined;
  if (tipoItem && VALID_TIPOS.includes(tipoItem as TipoItem) && itemId > 0) {
    initialItem = {
      ...emptyMovimientoItem(),
      tipo_item: tipoItem as TipoItem,
      item_id: itemId,
      item_nombre: nombre,
    };
  }

  return (
    <div className="max-w-3xl">
      <PageHeader
        title="Traslado bodega → vitrina"
        description="Mueve stock de bodega a vitrina."
        backHref="/inventario/movimientos"
        backLabel="Volver a movimientos"
      />
      <TrasladoForm initialItem={initialItem} />
    </div>
  );
}

export default function TrasladoPage() {
  return (
    <Suspense>
      <TrasladoView />
    </Suspense>
  );
}
