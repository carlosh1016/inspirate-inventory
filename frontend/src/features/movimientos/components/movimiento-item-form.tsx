'use client';

import { X } from 'lucide-react';

import { Combobox } from '@/components/forms/combobox';
import { DecimalInput } from '@/components/forms/decimal-input';
import { SelectField } from '@/components/forms/select-field';
import { FormError } from '@/components/forms/form-error';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { TipoItem, Ubicacion } from '@/types/domain';
import { searchCatalogItems, type CatalogItemOption } from '../api/search-items';

export interface MovimientoItemInput {
  tipo_item: '' | TipoItem;
  item_id: number;
  item_nombre: string;
  ubicacion: Ubicacion;
  cantidad: string;
  motivo: string;
}

export function emptyMovimientoItem(): MovimientoItemInput {
  return { tipo_item: '', item_id: 0, item_nombre: '', ubicacion: 'vitrina', cantidad: '', motivo: '' };
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

export type ItemErrors = Partial<Record<'tipo_item' | 'item_id' | 'cantidad' | 'motivo', string>>;

interface Props {
  index: number;
  value: MovimientoItemInput;
  onChange: (item: MovimientoItemInput) => void;
  onRemove?: () => void;
  showUbicacion?: boolean;
  showMotivo?: boolean;
  errors?: ItemErrors;
  /** e.g. stock-insuficiente message mapped to this row. */
  externalError?: string;
}

export function MovimientoItemForm({
  index,
  value,
  onChange,
  onRemove,
  showUbicacion = false,
  showMotivo = false,
  errors,
  externalError,
}: Props) {
  const isFragancia = value.tipo_item === 'fragancia';

  return (
    <div className="space-y-3 rounded-lg border border-border p-4">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">Ítem {index + 1}</span>
        {onRemove && (
          <Button type="button" variant="ghost" size="icon-sm" onClick={onRemove} aria-label="Quitar ítem">
            <X className="size-4" />
          </Button>
        )}
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <SelectField
          label="Tipo"
          value={value.tipo_item}
          onChange={(v) =>
            onChange({ ...value, tipo_item: v as TipoItem, item_id: 0, item_nombre: '' })
          }
          options={TIPO_ITEM_OPTIONS}
          placeholder="Tipo de ítem…"
          error={errors?.tipo_item}
        />

        <div className="space-y-2">
          <Label>Ítem</Label>
          <Combobox<CatalogItemOption>
            value={value.item_id > 0 ? value.item_id : null}
            selectedLabel={value.item_nombre || null}
            onChange={(id, label) =>
              onChange({ ...value, item_id: id ?? 0, item_nombre: label ?? '' })
            }
            searchFn={(q) => searchCatalogItems(value.tipo_item as TipoItem, q)}
            disabled={!value.tipo_item}
            placeholder={value.tipo_item ? 'Buscar ítem…' : 'Elige un tipo primero'}
            renderOption={(o) => (
              <span>
                <span className="font-medium">{o.label}</span>
                <span className="ml-2 text-xs text-muted-foreground">{o.detail}</span>
              </span>
            )}
          />
          <FormError message={errors?.item_id} />
        </div>

        {showUbicacion && (
          <SelectField
            label="Ubicación"
            value={value.ubicacion}
            onChange={(v) => onChange({ ...value, ubicacion: v as Ubicacion })}
            options={UBICACION_OPTIONS}
          />
        )}

        <div className="space-y-2">
          <Label htmlFor={`cantidad-${index}`}>Cantidad</Label>
          <DecimalInput
            id={`cantidad-${index}`}
            value={value.cantidad}
            onChange={(v) => onChange({ ...value, cantidad: v })}
            suffix={isFragancia ? 'g' : undefined}
          />
          <FormError message={errors?.cantidad} />
        </div>
      </div>

      {showMotivo && (
        <div className="space-y-2">
          <Label htmlFor={`motivo-${index}`}>Motivo</Label>
          <Input
            id={`motivo-${index}`}
            value={value.motivo}
            onChange={(e) => onChange({ ...value, motivo: e.target.value })}
          />
          <FormError message={errors?.motivo} />
        </div>
      )}

      {externalError && (
        <p className="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">{externalError}</p>
      )}
    </div>
  );
}
