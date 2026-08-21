'use client';

import { Combobox } from '@/components/forms/combobox';
import { SelectFilter } from '@/components/filters/select-filter';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { usePermission } from '@/hooks/use-permission';
import type { TipoItem } from '@/types/domain';
import { useModelosFullLookup } from '@/features/modelos-envase/api/use-modelos-full-lookup';
import { searchCatalogItems, type CatalogItemOption } from '../api/search-items';
import { useUsuariosSelect } from '../api/use-usuarios-select';
import type { MovimientosFilters } from '../api/use-movimientos';

interface Props {
  filters: MovimientosFilters;
  setFilter: <K extends keyof MovimientosFilters>(key: K, value: MovimientosFilters[K]) => void;
  setFilters: (partial: Partial<MovimientosFilters>) => void;
}

const TIPO_OPTIONS = [
  { value: 'all', label: 'Todos los tipos' },
  { value: 'entrada_mercancia', label: 'Entrada' },
  { value: 'traslado_bodega_vitrina', label: 'Traslado' },
  { value: 'venta', label: 'Venta' },
  { value: 'ajuste', label: 'Ajuste' },
  { value: 'danado', label: 'Dañado' },
  { value: 'devolucion', label: 'Devolución' },
  { value: 'correccion', label: 'Corrección' },
];

const TIPO_ITEM_OPTIONS = [
  { value: 'all', label: 'Todos los ítems' },
  { value: 'fragancia', label: 'Fragancias' },
  { value: 'variante_envase', label: 'Envases' },
  { value: 'producto', label: 'Productos' },
];

const UBICACION_OPTIONS = [
  { value: 'all', label: 'Toda ubicación' },
  { value: 'vitrina', label: 'Vitrina' },
  { value: 'bodega', label: 'Bodega' },
];

export function MovimientosFiltersBar({ filters, setFilter, setFilters }: Props) {
  const { isAdmin } = usePermission();
  const { data: usuarios } = useUsuariosSelect(isAdmin);
  const { data: modelosMap } = useModelosFullLookup();

  const usuarioOptions = [
    { value: '0', label: 'Todos los usuarios' },
    ...(usuarios ?? []).map((u) => ({ value: String(u.id), label: u.nombre_completo })),
  ];

  const isSpecificTipoItem = filters.tipo_item !== 'all';

  return (
    <div className="mb-4 flex flex-wrap items-end gap-3">
      <SelectFilter
        value={filters.tipo}
        onChange={(v) => setFilter('tipo', v)}
        options={TIPO_OPTIONS}
        ariaLabel="Filtrar por tipo de movimiento"
      />
      <SelectFilter
        value={filters.tipo_item}
        onChange={(v) => setFilters({ tipo_item: v, item_id: 0, item_nombre: '' })}
        options={TIPO_ITEM_OPTIONS}
        ariaLabel="Filtrar por tipo de ítem"
      />
      {isSpecificTipoItem && (
        <div className="w-56">
          <Combobox<CatalogItemOption>
            value={filters.item_id > 0 ? filters.item_id : null}
            selectedLabel={filters.item_nombre || null}
            onChange={(id, label) =>
              setFilters({ item_id: id ?? 0, item_nombre: label ?? '' })
            }
            searchFn={(q) => searchCatalogItems(filters.tipo_item as TipoItem, q, modelosMap)}
            placeholder="Ítem específico…"
            renderOption={(o) => (
              <span>
                <span className="font-medium">{o.label}</span>
                <span className="ml-2 text-xs text-muted-foreground">{o.detail}</span>
              </span>
            )}
          />
        </div>
      )}
      <SelectFilter
        value={filters.ubicacion}
        onChange={(v) => setFilter('ubicacion', v)}
        options={UBICACION_OPTIONS}
        ariaLabel="Filtrar por ubicación"
      />
      {isAdmin && (
        <SelectFilter
          value={String(filters.usuario_id)}
          onChange={(v) => setFilter('usuario_id', Number(v))}
          options={usuarioOptions}
          ariaLabel="Filtrar por usuario"
        />
      )}
      <div className="space-y-1">
        <Label htmlFor="fecha_desde" className="text-xs text-muted-foreground">
          Desde
        </Label>
        <Input
          id="fecha_desde"
          type="date"
          value={filters.fecha_desde}
          onChange={(e) => setFilter('fecha_desde', e.target.value)}
          className="w-40"
        />
      </div>
      <div className="space-y-1">
        <Label htmlFor="fecha_hasta" className="text-xs text-muted-foreground">
          Hasta
        </Label>
        <Input
          id="fecha_hasta"
          type="date"
          value={filters.fecha_hasta}
          onChange={(e) => setFilter('fecha_hasta', e.target.value)}
          className="w-40"
        />
      </div>
    </div>
  );
}
