'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Info, Plus } from 'lucide-react';
import { toast } from 'sonner';

import { Button, buttonVariants } from '@/components/ui/button';
import { getErrorMessage } from '@/lib/errors';
import type { TipoItem } from '@/types/domain';
import { useTraslado } from '../api/use-traslado';
import { parseStockInsuficiente } from '../stock-insuficiente';
import {
  MovimientoItemForm,
  emptyMovimientoItem,
  type ItemErrors,
  type MovimientoItemInput,
} from './movimiento-item-form';

function validateItem(item: MovimientoItemInput): ItemErrors {
  const errors: ItemErrors = {};
  if (!item.tipo_item) errors.tipo_item = 'Selecciona un tipo';
  if (item.item_id <= 0) errors.item_id = 'Selecciona un ítem';
  if (!/^\d+(\.\d+)?$/.test(item.cantidad) || Number.parseFloat(item.cantidad) <= 0) {
    errors.cantidad = 'Cantidad inválida';
  }
  return errors;
}

export function TrasladoForm({ initialItem }: { initialItem?: MovimientoItemInput }) {
  const router = useRouter();
  const traslado = useTraslado();
  const [items, setItems] = useState<MovimientoItemInput[]>([
    initialItem ?? emptyMovimientoItem(),
  ]);
  const [errors, setErrors] = useState<ItemErrors[]>([]);
  const [externalErrors, setExternalErrors] = useState<(string | undefined)[]>([]);

  const updateItem = (index: number, next: MovimientoItemInput) =>
    setItems((prev) => prev.map((it, i) => (i === index ? next : it)));

  const removeItem = (index: number) =>
    setItems((prev) => (prev.length === 1 ? prev : prev.filter((_, i) => i !== index)));

  const onSubmit = async () => {
    const itemErrors = items.map(validateItem);
    setErrors(itemErrors);
    setExternalErrors([]);
    if (itemErrors.some((e) => Object.keys(e).length > 0)) return;

    try {
      const payload = items.map((it) => ({
        tipo_item: it.tipo_item as TipoItem,
        item_id: it.item_id,
        cantidad: it.cantidad,
      }));
      await traslado.mutateAsync(payload);
      toast.success(`Se trasladaron ${items.length} ${items.length === 1 ? 'ítem' : 'ítems'}`);
      router.push('/inventario/movimientos');
    } catch (err) {
      const insuficientes = parseStockInsuficiente(err);
      if (insuficientes) {
        setExternalErrors(
          items.map((it) => {
            const match = insuficientes.find(
              (s) => s.tipo_item === it.tipo_item && s.item_id === it.item_id,
            );
            return match
              ? `Bodega disponible ${match.disponible}, requerido ${match.requerido}`
              : undefined;
          }),
        );
      } else {
        toast.error(getErrorMessage(err));
      }
    }
  };

  return (
    <div className="space-y-4">
      <p className="flex items-center gap-2 rounded-md bg-info/10 px-3 py-2 text-sm text-info">
        <Info className="size-4 shrink-0" />
        Los ítems seleccionados se moverán de bodega a vitrina.
      </p>

      {items.map((item, index) => (
        <MovimientoItemForm
          key={index}
          index={index}
          value={item}
          onChange={(next) => updateItem(index, next)}
          onRemove={items.length > 1 ? () => removeItem(index) : undefined}
          errors={errors[index]}
          externalError={externalErrors[index]}
        />
      ))}

      <Button
        type="button"
        variant="outline"
        onClick={() => setItems((prev) => [...prev, emptyMovimientoItem()])}
      >
        <Plus className="size-4" />
        Agregar otro ítem
      </Button>

      <div className="flex justify-end gap-2 border-t border-border pt-4">
        <Link href="/inventario/movimientos" className={buttonVariants({ variant: 'outline' })}>
          Cancelar
        </Link>
        <Button type="button" onClick={onSubmit} disabled={traslado.isPending}>
          {traslado.isPending ? 'Trasladando…' : 'Registrar traslado'}
        </Button>
      </div>
    </div>
  );
}
