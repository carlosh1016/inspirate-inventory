'use client';

import { BooleanFilter } from '@/components/filters/boolean-filter';
import { SelectFilter } from '@/components/filters/select-filter';
import { MoneyInput } from '@/components/forms/money-input';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { usePermission } from '@/hooks/use-permission';
import { useMetodosPagoActivos } from '../api/use-metodos-pago-activos';
import { useUsuariosSelect } from '@/features/movimientos/api/use-usuarios-select';
import type { VentasFilters } from '../types';

interface Props {
  filters: VentasFilters;
  setFilter: <K extends keyof VentasFilters>(key: K, value: VentasFilters[K]) => void;
}

export function VentasFiltersBar({ filters, setFilter }: Props) {
  const { isAdmin } = usePermission();
  const { data: metodos } = useMetodosPagoActivos();
  const { data: usuarios } = useUsuariosSelect(isAdmin);

  const metodoOptions = [
    { value: '0', label: 'Todos los métodos' },
    ...(metodos ?? []).map((m) => ({ value: String(m.id), label: m.nombre })),
  ];
  const usuarioOptions = [
    { value: '0', label: 'Todas las vendedoras' },
    ...(usuarios ?? []).map((u) => ({ value: String(u.id), label: u.nombre_completo })),
  ];

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
        value={String(filters.metodo_pago_id)}
        onChange={(v) => setFilter('metodo_pago_id', Number(v))}
        options={metodoOptions}
        ariaLabel="Filtrar por método de pago"
      />
      {isAdmin && (
        <SelectFilter
          value={String(filters.usuario_id)}
          onChange={(v) => setFilter('usuario_id', Number(v))}
          options={usuarioOptions}
          ariaLabel="Filtrar por vendedora"
        />
      )}
      <div className="space-y-1">
        <Label className="text-xs text-muted-foreground">Total mín.</Label>
        <MoneyInput value={filters.total_min} onChange={(v) => setFilter('total_min', v)} className="w-32" />
      </div>
      <div className="space-y-1">
        <Label className="text-xs text-muted-foreground">Total máx.</Label>
        <MoneyInput value={filters.total_max} onChange={(v) => setFilter('total_max', v)} className="w-32" />
      </div>
      <BooleanFilter
        id="con_descuento"
        label="Con descuento"
        checked={filters.con_descuento}
        onChange={(v) => setFilter('con_descuento', v)}
      />
    </div>
  );
}
