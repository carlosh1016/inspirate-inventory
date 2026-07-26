'use client';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { formatCurrency, formatDate } from '@/lib/formatters';
import { cn } from '@/lib/utils';
import type { Meta } from '@/types/api';
import type { CuadreListItem } from '../types';

interface Props {
  data: CuadreListItem[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  onPageChange: (page: number) => void;
  onRowClick: (row: CuadreListItem) => void;
}

const columns: Column<CuadreListItem>[] = [
  { key: 'fecha', header: 'Fecha', cell: (c) => formatDate(c.fecha) },
  {
    key: 'estado',
    header: 'Estado',
    cell: (c) => (
      <span
        className={cn(
          'rounded-md px-2 py-0.5 text-xs font-medium',
          c.estado === 'abierto' ? 'bg-success/10 text-success' : 'bg-muted text-muted-foreground',
        )}
      >
        {c.estado === 'abierto' ? 'Abierta' : 'Cerrada'}
      </span>
    ),
  },
  {
    key: 'total_ventas',
    header: 'Total ventas',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (c) => formatCurrency(c.total_ventas),
  },
  {
    key: 'saldo',
    header: 'Saldo',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (c) => formatCurrency(c.saldo_calculado),
  },
  {
    key: 'cerrado_por',
    header: 'Cerrado por',
    cell: (c) => c.cerrado_por?.nombre_completo ?? <span className="text-muted-foreground">—</span>,
  },
];

export function CuadresTable({ data, meta, isLoading, page, onPageChange, onRowClick }: Props) {
  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      emptyMessage="No hay cuadres que coincidan con los filtros."
      rowKey={(c) => c.id}
      onRowClick={onRowClick}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
