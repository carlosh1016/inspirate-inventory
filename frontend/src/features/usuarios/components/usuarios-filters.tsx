'use client';

import { SelectFilter } from '@/components/filters/select-filter';
import type { UsuariosFilters } from '../types';

interface Props {
  filters: UsuariosFilters;
  setFilter: <K extends keyof UsuariosFilters>(key: K, value: UsuariosFilters[K]) => void;
}

const ROL_OPTIONS = [
  { value: 'all', label: 'Todos los roles' },
  { value: 'admin', label: 'Admin' },
  { value: 'vendedora', label: 'Vendedora' },
];

const ACTIVO_OPTIONS = [
  { value: 'all', label: 'Todos los estados' },
  { value: 'true', label: 'Activos' },
  { value: 'false', label: 'Inactivos' },
];

export function UsuariosFiltersBar({ filters, setFilter }: Props) {
  return (
    <div className="mb-4 flex flex-wrap items-end gap-3">
      <SelectFilter
        value={filters.rol}
        onChange={(v) => setFilter('rol', v as UsuariosFilters['rol'])}
        options={ROL_OPTIONS}
        ariaLabel="Filtrar por rol"
      />
      <SelectFilter
        value={filters.activo}
        onChange={(v) => setFilter('activo', v as UsuariosFilters['activo'])}
        options={ACTIVO_OPTIONS}
        ariaLabel="Filtrar por estado"
      />
    </div>
  );
}
