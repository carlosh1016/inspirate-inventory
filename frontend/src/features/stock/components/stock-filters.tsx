'use client';

import { BooleanFilter } from '@/components/filters/boolean-filter';
import { SelectFilter } from '@/components/filters/select-filter';
import type { StockFilters } from '../types';

interface Props {
  filters: StockFilters;
  setFilter: <K extends keyof StockFilters>(key: K, value: StockFilters[K]) => void;
}

const TIPO_OPTIONS = [
  { value: 'all', label: 'Todos los tipos' },
  { value: 'fragancia', label: 'Fragancias' },
  { value: 'variante_envase', label: 'Envases' },
  { value: 'producto', label: 'Productos' },
];

const UBICACION_OPTIONS = [
  { value: 'all', label: 'Vitrina + bodega' },
  { value: 'vitrina', label: 'Solo vitrina' },
  { value: 'bodega', label: 'Solo bodega' },
];

export function StockFiltersBar({ filters, setFilter }: Props) {
  return (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      <SelectFilter
        value={filters.tipo_item}
        onChange={(v) => setFilter('tipo_item', v)}
        options={TIPO_OPTIONS}
        ariaLabel="Filtrar por tipo"
      />
      <SelectFilter
        value={filters.ubicacion}
        onChange={(v) => setFilter('ubicacion', v)}
        options={UBICACION_OPTIONS}
        ariaLabel="Filtrar por ubicación"
      />
      <BooleanFilter
        id="stock_bajo"
        label="Bajo stock"
        checked={filters.stock_bajo}
        onChange={(v) => setFilter('stock_bajo', v)}
      />
      <BooleanFilter
        id="stock_cero"
        label="Stock cero"
        checked={filters.stock_cero}
        onChange={(v) => setFilter('stock_cero', v)}
      />
    </div>
  );
}
