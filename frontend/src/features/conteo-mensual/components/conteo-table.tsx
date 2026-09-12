'use client';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { cn } from '@/lib/utils';
import type { Meta } from '@/types/api';
import { formatPeriodo } from '../format-periodo';
import type { ConteoListItem } from '../types';

interface Props {
  data: ConteoListItem[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  onPageChange: (page: number) => void;
  onRowClick: (row: ConteoListItem) => void;
}

const columns: Column<ConteoListItem>[] = [
  { key: 'periodo', header: 'Periodo', cell: (c) => formatPeriodo(c.periodo) },
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
        {c.estado === 'abierto' ? 'Abierto' : 'Cerrado'}
      </span>
    ),
  },
  { key: 'creado_por', header: 'Creado por', cell: (c) => c.creado_por.nombre_completo },
  {
    key: 'cerrado_por',
    header: 'Cerrado por',
    cell: (c) => c.cerrado_por?.nombre_completo ?? <span className="text-muted-foreground">—</span>,
  },
];

export function ConteoTable({ data, meta, isLoading, page, onPageChange, onRowClick }: Props) {
  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      emptyMessage="Todavía no se ha creado ningún conteo mensual."
      rowKey={(c) => c.id}
      onRowClick={onRowClick}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
