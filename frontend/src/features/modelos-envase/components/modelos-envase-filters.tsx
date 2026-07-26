'use client';

import { SearchInput } from '@/components/filters/search-input';
import { SelectFilter } from '@/components/filters/select-filter';
import type { ModelosEnvaseFilters } from '../types';

interface Props {
  filters: ModelosEnvaseFilters;
  setFilter: <K extends keyof ModelosEnvaseFilters>(
    key: K,
    value: ModelosEnvaseFilters[K],
  ) => void;
}

const ESTADO_OPTIONS = [
  { value: 'true', label: 'Activos' },
  { value: 'false', label: 'Inactivos' },
  { value: 'all', label: 'Todos' },
];

export function ModelosEnvaseFiltersBar({ filters, setFilter }: Props) {
  return (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      <SearchInput
        value={filters.q}
        onChange={(v) => setFilter('q', v)}
        placeholder="Buscar por tipo…"
        className="w-full sm:w-64"
      />
      <SelectFilter
        value={filters.activo}
        onChange={(v) => setFilter('activo', v)}
        options={ESTADO_OPTIONS}
        ariaLabel="Filtrar por estado"
      />
    </div>
  );
}
