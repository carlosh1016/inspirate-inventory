'use client';

import { BooleanFilter } from '@/components/filters/boolean-filter';
import { SearchInput } from '@/components/filters/search-input';
import { SelectFilter } from '@/components/filters/select-filter';
import { CATEGORIA_OPTIONS, type ProductosFilters } from '../types';

interface Props {
  filters: ProductosFilters;
  setFilter: <K extends keyof ProductosFilters>(key: K, value: ProductosFilters[K]) => void;
}

const ESTADO_OPTIONS = [
  { value: 'true', label: 'Activos' },
  { value: 'false', label: 'Inactivos' },
  { value: 'all', label: 'Todos' },
];

const CATEGORIA_FILTER_OPTIONS = [
  { value: 'all', label: 'Todas las categorías' },
  ...CATEGORIA_OPTIONS,
];

export function ProductosFiltersBar({ filters, setFilter }: Props) {
  return (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      <SearchInput
        value={filters.q}
        onChange={(v) => setFilter('q', v)}
        placeholder="Buscar por nombre…"
        className="w-full sm:w-56"
      />
      <SelectFilter
        value={filters.categoria}
        onChange={(v) => setFilter('categoria', v)}
        options={CATEGORIA_FILTER_OPTIONS}
        ariaLabel="Filtrar por categoría"
      />
      <SelectFilter
        value={filters.activo}
        onChange={(v) => setFilter('activo', v)}
        options={ESTADO_OPTIONS}
        ariaLabel="Filtrar por estado"
      />
      <BooleanFilter
        id="stock_bajo"
        label="Bajo stock"
        checked={filters.stock_bajo}
        onChange={(v) => setFilter('stock_bajo', v)}
      />
    </div>
  );
}
