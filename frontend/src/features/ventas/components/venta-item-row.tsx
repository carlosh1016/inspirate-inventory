'use client';

import { X } from 'lucide-react';

import { Combobox } from '@/components/forms/combobox';
import { DecimalInput } from '@/components/forms/decimal-input';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { formatCurrency } from '@/lib/formatters';
import type { ModeloEnvase } from '@/features/modelos-envase/types';
import {
  searchEnvasesVenta,
  searchFeromonas,
  searchFraganciasVenta,
  searchProductosVenta,
  type EnvaseOption,
  type FraganciaOption,
  type ProductoOption,
} from '../api/search-venta-items';
import type { TipoLinea, VentaItemState } from '../types';

const TIPO_LABEL: Record<TipoLinea, string> = {
  envase_con_fragancia: 'Envase con fragancia',
  recarga: 'Recarga',
  envase_solo: 'Envase solo',
  producto_otro: 'Otro',
};

interface Props {
  index: number;
  value: VentaItemState;
  onChange: (item: VentaItemState) => void;
  onRemove: () => void;
  modelosMap: Map<number, ModeloEnvase>;
  subtotal: number;
  error?: string;
}

function renderOptionWithDetail(o: { label: string; detail: string }) {
  return (
    <span>
      <span className="font-medium">{o.label}</span>
      {o.detail && <span className="ml-2 text-xs text-muted-foreground">{o.detail}</span>}
    </span>
  );
}

export function VentaItemRow({ index, value, onChange, onRemove, modelosMap, subtotal, error }: Props) {
  const esFragancia = value.tipo_linea === 'envase_con_fragancia' || value.tipo_linea === 'recarga';
  const usaEnvase = esFragancia || value.tipo_linea === 'envase_solo';

  return (
    <div className="space-y-3 rounded-lg border border-border p-4">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium">{TIPO_LABEL[value.tipo_linea]}</span>
        <Button type="button" variant="ghost" size="icon-sm" onClick={onRemove} aria-label="Quitar ítem">
          <X className="size-4" />
        </Button>
      </div>

      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        {usaEnvase && (
          <div className="space-y-2">
            <Label>Envase</Label>
            <Combobox<EnvaseOption>
              value={value.envase?.variante_id ?? null}
              selectedLabel={value.envase?.label ?? null}
              onChange={(id, _label, option) =>
                onChange({
                  ...value,
                  envase: option
                    ? {
                        variante_id: option.id,
                        label: option.label,
                        equiv_gramos: option.equiv_gramos,
                        precio_solo: option.precio_solo,
                        precio_con_fragancia: option.precio_con_fragancia,
                        precio_recarga: option.precio_recarga,
                      }
                    : null,
                  // Auto-llenar gramos con la equivalencia del envase (editable).
                  gramos: esFragancia && option ? option.equiv_gramos : value.gramos,
                })
              }
              searchFn={(q) => searchEnvasesVenta(q, modelosMap)}
              placeholder="Buscar envase…"
              renderOption={renderOptionWithDetail}
            />
          </div>
        )}

        {esFragancia && (
          <div className="space-y-2">
            <Label>Fragancia</Label>
            <Combobox<FraganciaOption>
              value={value.fragancia?.id ?? null}
              selectedLabel={value.fragancia?.label ?? null}
              onChange={(id, label) =>
                onChange({ ...value, fragancia: id ? { id, label: label ?? '' } : null })
              }
              searchFn={searchFraganciasVenta}
              placeholder="Buscar fragancia…"
              renderOption={renderOptionWithDetail}
            />
          </div>
        )}

        {value.tipo_linea === 'producto_otro' && (
          <div className="space-y-2">
            <Label>Producto</Label>
            <Combobox<ProductoOption>
              value={value.producto?.id ?? null}
              selectedLabel={value.producto?.label ?? null}
              onChange={(id, label, option) =>
                onChange({
                  ...value,
                  producto: option ? { id: option.id, label: option.label, precio: option.precio } : null,
                })
              }
              searchFn={searchProductosVenta}
              placeholder="Buscar producto…"
              renderOption={renderOptionWithDetail}
            />
          </div>
        )}

        {esFragancia && (
          <div className="space-y-2">
            <Label htmlFor={`gramos-${index}`}>Gramos</Label>
            <DecimalInput
              id={`gramos-${index}`}
              suffix="g"
              value={value.gramos}
              onChange={(gramos) => onChange({ ...value, gramos })}
            />
          </div>
        )}

        <div className="space-y-2">
          <Label htmlFor={`cantidad-${index}`}>Cantidad</Label>
          <Input
            id={`cantidad-${index}`}
            type="number"
            min={1}
            value={value.cantidad}
            onChange={(e) => {
              const n = Number.parseInt(e.target.value, 10);
              onChange({ ...value, cantidad: Number.isNaN(n) ? 0 : n });
            }}
          />
        </div>
      </div>

      {esFragancia && (
        <div className="space-y-2">
          <div className="flex items-center gap-2">
            <Checkbox
              id={`feromona-${index}`}
              checked={value.feromona_enabled}
              onCheckedChange={(checked) =>
                onChange({
                  ...value,
                  feromona_enabled: checked === true,
                  feromona: checked === true ? value.feromona : null,
                })
              }
            />
            <Label htmlFor={`feromona-${index}`} className="cursor-pointer text-sm font-normal">
              Agregar feromonas
            </Label>
          </div>
          {value.feromona_enabled && (
            <Combobox<ProductoOption>
              value={value.feromona?.id ?? null}
              selectedLabel={value.feromona?.label ?? null}
              onChange={(id, label, option) =>
                onChange({
                  ...value,
                  feromona: option ? { id: option.id, label: option.label, precio: option.precio } : null,
                })
              }
              searchFn={searchFeromonas}
              placeholder="Buscar feromona…"
              renderOption={(o) => (
                <span>
                  <span className="font-medium">{o.label}</span>
                  <span className="ml-2 text-xs text-muted-foreground">{formatCurrency(o.precio)}</span>
                </span>
              )}
            />
          )}
        </div>
      )}

      <div className="flex items-center justify-between border-t border-border pt-2">
        {error ? (
          <p className="text-sm text-destructive">{error}</p>
        ) : (
          <span />
        )}
        <span className="text-sm font-medium tabular-nums">Subtotal: {formatCurrency(subtotal)}</span>
      </div>
    </div>
  );
}
