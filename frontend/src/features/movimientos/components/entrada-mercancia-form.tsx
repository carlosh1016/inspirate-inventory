'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { Plus } from 'lucide-react';
import { toast } from 'sonner';

import { Button, buttonVariants } from '@/components/ui/button';
import { getErrorMessage } from '@/lib/errors';
import type { TipoItem } from '@/types/domain';
import { useEntradaMercancia } from '../api/use-entrada-mercancia';
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

export function EntradaMercanciaForm() {
  const router = useRouter();
  const entrada = useEntradaMercancia();
  const [items, setItems] = useState<MovimientoItemInput[]>([emptyMovimientoItem()]);
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
        ubicacion: it.ubicacion,
        cantidad: it.cantidad,
        motivo: it.motivo || undefined,
      }));
      await entrada.mutateAsync(payload);
      toast.success(`Se registraron ${items.length} ${items.length === 1 ? 'ítem' : 'ítems'}`);
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
              ? `Disponible ${match.disponible}, requerido ${match.requerido}`
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
      {items.map((item, index) => (
        <MovimientoItemForm
          key={index}
          index={index}
          value={item}
          onChange={(next) => updateItem(index, next)}
          onRemove={items.length > 1 ? () => removeItem(index) : undefined}
          showUbicacion
          showMotivo
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
        <Button type="button" onClick={onSubmit} disabled={entrada.isPending}>
          {entrada.isPending ? 'Registrando…' : 'Registrar entrada'}
        </Button>
      </div>
    </div>
  );
}
