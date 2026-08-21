'use client';

import { BooleanFilter } from '@/components/filters/boolean-filter';
import { SearchInput } from '@/components/filters/search-input';
import { SelectFilter } from '@/components/filters/select-filter';
import { useModelosLookup } from '@/features/modelos-envase/api/use-modelos-lookup';
import type { VariantesEnvaseFilters } from '../types';

interface Props {
  filters: VariantesEnvaseFilters;
  setFilter: <K extends keyof VariantesEnvaseFilters>(
    key: K,
    value: VariantesEnvaseFilters[K],
  ) => void;
}

const ESTADO_OPTIONS = [
  { value: 'true', label: 'Activas' },
  { value: 'false', label: 'Inactivas' },
  { value: 'all', label: 'Todas' },
];

export function VariantesEnvaseFiltersBar({ filters, setFilter }: Props) {
  const { data: modelosMap } = useModelosLookup();

  const modeloOptions = [
    { value: '0', label: 'Todos los modelos' },
    ...[...(modelosMap?.entries() ?? [])].map(([id, label]) => ({ value: String(id), label })),
  ];

  return (
    <div className="mb-4 flex flex-wrap items-center gap-3">
      <SearchInput
        value={filters.q}
        onChange={(v) => setFilter('q', v)}
        placeholder="Buscar por grosor…"
        className="w-full sm:w-56"
      />
      <SelectFilter
        value={String(filters.modelo_envase_id)}
        onChange={(v) => setFilter('modelo_envase_id', Number(v))}
        options={modeloOptions}
        ariaLabel="Filtrar por modelo"
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
