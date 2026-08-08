'use client';

import type { ReactNode } from 'react';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { StockBadge } from '@/components/stock-badge';
import { formatCodigoFragancia, formatGramos } from '@/lib/formatters';
import type { Meta } from '@/types/api';
import type { Fragancia } from '../types';

interface Props {
  data: Fragancia[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  onPageChange: (page: number) => void;
  onRowClick: (row: Fragancia) => void;
  emptyAction?: ReactNode;
}

const columns: Column<Fragancia>[] = [
  {
    key: 'numero_genero',
    header: 'Código',
    cell: (f) => (
      <span className="font-mono text-xs text-muted-foreground">
        {formatCodigoFragancia(f.genero, f.numero_genero)}
      </span>
    ),
  },
  {
    key: 'nombre_comercial',
    header: 'Nombre comercial',
    cell: (f) => (
      <span className="font-medium">
        {f.nombre_comercial}
        {!f.activo && <span className="ml-2 text-xs text-muted-foreground">(inactiva)</span>}
      </span>
    ),
  },
  {
    key: 'nombre_alternativo',
    header: 'Alternativo',
    cell: (f) => f.nombre_alternativo ?? <span className="text-muted-foreground">—</span>,
  },
  {
    key: 'genero',
    header: 'Género',
    cell: (f) => <span className="capitalize">{f.genero}</span>,
  },
  {
    key: 'vitrina',
    header: 'Vitrina',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (f) => formatGramos(f.stock.vitrina),
  },
  {
    key: 'bodega',
    header: 'Bodega',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (f) => formatGramos(f.stock.bodega),
  },
  {
    key: 'total',
    header: 'Total',
    className: 'text-right',
    headerClassName: 'text-right',
    cell: (f) => <StockBadge total={f.stock.total} minimo={f.gramos_minimo} unidad="gramos" />,
  },
];

export function FraganciasTable({
  data,
  meta,
  isLoading,
  page,
  onPageChange,
  onRowClick,
  emptyAction,
}: Props) {
  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      emptyMessage="No hay fragancias que coincidan con los filtros."
      emptyAction={emptyAction}
      rowKey={(f) => f.id}
      onRowClick={onRowClick}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
