'use client';

import type { ReactNode } from 'react';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { formatCurrency } from '@/lib/formatters';
import type { Meta } from '@/types/api';
import type { ModeloEnvase } from '../types';

interface Props {
  data: ModeloEnvase[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  onPageChange: (page: number) => void;
  onRowClick: (row: ModeloEnvase) => void;
  emptyAction?: ReactNode;
}

const columns: Column<ModeloEnvase>[] = [
  {
    key: 'tipo',
    header: 'Tipo',
    cell: (m) => (
      <span className="font-medium">
        {m.tipo}
        {!m.activo && <span className="ml-2 text-xs text-muted-foreground">(inactivo)</span>}
      </span>
    ),
  },
  { key: 'tamano_oz', header: 'Tamaño', className: 'tabular-nums', cell: (m) => `${m.tamano_oz} oz` },
  { key: 'equiv_gramos', header: 'Equiv.', className: 'tabular-nums', cell: (m) => `${m.equiv_gramos} g` },
  {
    key: 'precio_solo',
    header: 'Solo',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (m) => formatCurrency(m.precio_solo),
  },
  {
    key: 'precio_con_fragancia',
    header: 'Con fragancia',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (m) => formatCurrency(m.precio_con_fragancia),
  },
  {
    key: 'precio_recarga',
    header: 'Recarga',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (m) => formatCurrency(m.precio_recarga),
  },
  {
    key: 'variantes_activas',
    header: 'Variantes',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (m) => m.variantes_activas,
  },
];

export function ModelosEnvaseTable({
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
      emptyMessage="No hay modelos que coincidan con los filtros."
      emptyAction={emptyAction}
      rowKey={(m) => m.id}
      onRowClick={onRowClick}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
