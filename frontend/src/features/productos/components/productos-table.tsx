'use client';

import type { ReactNode } from 'react';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { StockBadge } from '@/components/stock-badge';
import { formatCurrency } from '@/lib/formatters';
import type { Meta } from '@/types/api';
import { categoriaLabel, type Producto } from '../types';

interface Props {
  data: Producto[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  onPageChange: (page: number) => void;
  onRowClick: (row: Producto) => void;
  emptyAction?: ReactNode;
}

const columns: Column<Producto>[] = [
  {
    key: 'nombre',
    header: 'Nombre',
    cell: (p) => (
      <span className="font-medium">
        {p.nombre}
        {!p.activo && <span className="ml-2 text-xs text-muted-foreground">(inactivo)</span>}
      </span>
    ),
  },
  { key: 'categoria', header: 'Categoría', cell: (p) => categoriaLabel(p.categoria) },
  {
    key: 'precio',
    header: 'Precio',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (p) => formatCurrency(p.precio),
  },
  {
    key: 'total',
    header: 'Stock',
    className: 'text-right',
    headerClassName: 'text-right',
    cell: (p) => <StockBadge total={p.stock.total} minimo={p.stock_minimo} unidad="unidades" />,
  },
  {
    key: 'stock_minimo',
    header: 'Mínimo',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (p) => p.stock_minimo,
  },
];

export function ProductosTable({
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
      emptyMessage="No hay productos que coincidan con los filtros."
      emptyAction={emptyAction}
      rowKey={(p) => p.id}
      onRowClick={onRowClick}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
