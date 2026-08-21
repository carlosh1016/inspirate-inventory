'use client';

import { useState } from 'react';
import Link from 'next/link';
import { toast } from 'sonner';

import { Combobox } from '@/components/forms/combobox';
import { DecimalInput } from '@/components/forms/decimal-input';
import { FormError } from '@/components/forms/form-error';
import { SelectField } from '@/components/forms/select-field';
import { TextareaField } from '@/components/forms/textarea-field';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { getErrorMessage } from '@/lib/errors';
import type { TipoItem, Ubicacion } from '@/types/domain';
import { useModelosFullLookup } from '@/features/modelos-envase/api/use-modelos-full-lookup';
import { searchCatalogItems, type CatalogItemOption } from '../api/search-items';

export interface SingleMovimientoValues {
  tipo_item: TipoItem;
  item_id: number;
  ubicacion: Ubicacion;
  cantidad: string;
  motivo: string;
}

interface Props {
  cantidadLabel: string;
  cantidadHint?: string;
  /** ajuste/corrección permiten 0 (cantidad final); dañado exige > 0. */
  allowZero?: boolean;
  submitLabel: string;
  submittingLabel: string;
  onSubmit: (values: SingleMovimientoValues) => Promise<void>;
}

const TIPO_ITEM_OPTIONS = [
  { value: 'fragancia', label: 'Fragancia' },
  { value: 'variante_envase', label: 'Variante de envase' },
  { value: 'producto', label: 'Producto' },
];

const UBICACION_OPTIONS = [
  { value: 'vitrina', label: 'Vitrina' },
  { value: 'bodega', label: 'Bodega' },
];

type Errors = Partial<Record<'tipo_item' | 'item_id' | 'cantidad' | 'motivo', string>>;

export function SingleMovimientoForm({
  cantidadLabel,
  cantidadHint,
  allowZero = false,
  submitLabel,
  submittingLabel,
  onSubmit,
}: Props) {
  const [tipoItem, setTipoItem] = useState<'' | TipoItem>('');
  const [itemId, setItemId] = useState(0);
  const [itemNombre, setItemNombre] = useState<string | null>(null);
  const [ubicacion, setUbicacion] = useState<Ubicacion>('vitrina');
  const [cantidad, setCantidad] = useState('');
  const [motivo, setMotivo] = useState('');
  const [errors, setErrors] = useState<Errors>({});
  const [submitting, setSubmitting] = useState(false);
  const { data: modelosMap } = useModelosFullLookup();

  const validate = (): Errors => {
    const e: Errors = {};
    if (!tipoItem) e.tipo_item = 'Selecciona un tipo';
    if (itemId <= 0) e.item_id = 'Selecciona un ítem';
    const n = Number.parseFloat(cantidad);
    if (!/^\d+(\.\d+)?$/.test(cantidad) || (allowZero ? n < 0 : n <= 0)) {
      e.cantidad = 'Cantidad inválida';
    }
    if (!motivo.trim()) e.motivo = 'El motivo es obligatorio';
    return e;
  };

  const handleSubmit = async () => {
    const e = validate();
    setErrors(e);
    if (Object.keys(e).length > 0) return;

    setSubmitting(true);
    try {
      await onSubmit({
        tipo_item: tipoItem as TipoItem,
        item_id: itemId,
        ubicacion,
        cantidad,
        motivo,
      });
    } catch (err) {
      toast.error(getErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Card>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <SelectField
            label="Tipo"
            value={tipoItem}
            onChange={(v) => {
              setTipoItem(v as TipoItem);
              setItemId(0);
              setItemNombre(null);
            }}
            options={TIPO_ITEM_OPTIONS}
            placeholder="Tipo de ítem…"
            error={errors.tipo_item}
          />
          <div className="space-y-2">
            <Label>Ítem</Label>
            <Combobox<CatalogItemOption>
              value={itemId > 0 ? itemId : null}
              selectedLabel={itemNombre}
              onChange={(id, label) => {
                setItemId(id ?? 0);
                setItemNombre(label);
              }}
              searchFn={(q) => searchCatalogItems(tipoItem as TipoItem, q, modelosMap)}
              disabled={!tipoItem}
              placeholder={tipoItem ? 'Buscar ítem…' : 'Elige un tipo primero'}
              renderOption={(o) => (
                <span>
                  <span className="font-medium">{o.label}</span>
                  <span className="ml-2 text-xs text-muted-foreground">{o.detail}</span>
                </span>
              )}
            />
            <FormError message={errors.item_id} />
          </div>
        </div>

        <SelectField
          label="Ubicación"
          value={ubicacion}
          onChange={(v) => setUbicacion(v as Ubicacion)}
          options={UBICACION_OPTIONS}
          className="sm:max-w-xs"
        />

        <div className="space-y-2">
          <Label htmlFor="cantidad">{cantidadLabel}</Label>
          <DecimalInput
            id="cantidad"
            value={cantidad}
            onChange={setCantidad}
            suffix={tipoItem === 'fragancia' ? 'g' : undefined}
            className="sm:max-w-xs"
          />
          {cantidadHint && <p className="text-xs text-muted-foreground">{cantidadHint}</p>}
          <FormError message={errors.cantidad} />
        </div>

        <TextareaField
          id="motivo"
          label="Motivo"
          value={motivo}
          onChange={(e) => setMotivo(e.target.value)}
          error={errors.motivo}
        />

        <div className="flex justify-end gap-2 border-t border-border pt-4">
          <Link href="/inventario/movimientos" className={buttonVariants({ variant: 'outline' })}>
            Cancelar
          </Link>
          <Button type="button" onClick={handleSubmit} disabled={submitting}>
            {submitting ? submittingLabel : submitLabel}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
