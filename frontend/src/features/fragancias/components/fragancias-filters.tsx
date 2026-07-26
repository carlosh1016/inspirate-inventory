'use client';

import { BooleanFilter } from '@/components/filters/boolean-filter';
import { SearchInput } from '@/components/filters/search-input';
import { SelectFilter } from '@/components/filters/select-filter';
import type { FraganciasFilters } from '../types';

interface Props {
  filters: FraganciasFilters;
  setFilter: <K extends keyof FraganciasFilters>(key: K, value: FraganciasFilters[K]) => void;
}

const GENERO_OPTIONS = [
  { value: 'all', label: 'Todos los géneros' },
  { value: 'femenina', label: 'Femenina' },
  { value: 'masculina', label: 'Masculina' },
];

const ESTADO_OPTIONS = [
  { value: 'true', label: 'Activas' },
  { value: 'false', label: 'Inactivas' },
  { value: 'all', label: 'Todas' },
];

export function FraganciasFiltersBar({ filters, setFilter }: Props) {
  return (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      <SearchInput
        value={filters.q}
        onChange={(v) => setFilter('q', v)}
        placeholder="Buscar por nombre…"
        className="w-full sm:w-64"
      />
      <SelectFilter
        value={filters.genero}
        onChange={(v) => setFilter('genero', v)}
        options={GENERO_OPTIONS}
        ariaLabel="Filtrar por género"
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
