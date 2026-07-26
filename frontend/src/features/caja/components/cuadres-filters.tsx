'use client';

import { SelectFilter } from '@/components/filters/select-filter';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { CuadresFilters } from '../types';

interface Props {
  filters: CuadresFilters;
  setFilter: <K extends keyof CuadresFilters>(key: K, value: CuadresFilters[K]) => void;
}

const ESTADO_OPTIONS = [
  { value: 'all', label: 'Todos' },
  { value: 'abierto', label: 'Abierta' },
  { value: 'cerrado', label: 'Cerrada' },
];

export function CuadresFiltersBar({ filters, setFilter }: Props) {
  return (
    <div className="mb-4 flex flex-wrap items-end gap-3">
      <div className="space-y-1">
        <Label htmlFor="fecha_desde" className="text-xs text-muted-foreground">Desde</Label>
        <Input id="fecha_desde" type="date" value={filters.fecha_desde} onChange={(e) => setFilter('fecha_desde', e.target.value)} className="w-40" />
      </div>
      <div className="space-y-1">
        <Label htmlFor="fecha_hasta" className="text-xs text-muted-foreground">Hasta</Label>
        <Input id="fecha_hasta" type="date" value={filters.fecha_hasta} onChange={(e) => setFilter('fecha_hasta', e.target.value)} className="w-40" />
      </div>
      <SelectFilter
        value={filters.estado}
        onChange={(v) => setFilter('estado', v)}
        options={ESTADO_OPTIONS}
        ariaLabel="Filtrar por estado"
      />
    </div>
  );
}
