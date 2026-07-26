'use client';

import type { ReactNode } from 'react';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { StockBadge } from '@/components/stock-badge';
import { formatGramos } from '@/lib/formatters';
import { cn } from '@/lib/utils';
import type { Meta } from '@/types/api';
import type { TipoItem } from '@/types/domain';
import type { StockItem } from '../types';

interface Props {
  data: StockItem[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  onPageChange: (page: number) => void;
  emptyMessage?: string;
  emptyAction?: ReactNode;
  renderActions?: (row: StockItem) => ReactNode;
}

const TIPO_ITEM_BADGE: Record<TipoItem, { label: string; className: string }> = {
  fragancia: { label: 'Fragancia', className: 'bg-primary/10 text-primary' },
  variante_envase: { label: 'Envase', className: 'bg-info/10 text-info' },
  producto: { label: 'Producto', className: 'bg-success/10 text-success' },
};

function fmtStock(value: string, unidad: string): string {
  return unidad === 'gramos' ? formatGramos(value) : value;
}

export function StockTable({
  data,
  meta,
  isLoading,
  page,
  onPageChange,
  emptyMessage = 'No hay ítems que coincidan con los filtros.',
  emptyAction,
  renderActions,
}: Props) {
  const columns: Column<StockItem>[] = [
    {
      key: 'tipo',
      header: 'Tipo',
      cell: (s) => {
        const badge = TIPO_ITEM_BADGE[s.tipo_item];
        return (
          <span className={cn('rounded-md px-2 py-0.5 text-xs font-medium', badge.className)}>
            {badge.label}
          </span>
        );
      },
    },
    { key: 'nombre', header: 'Nombre', cell: (s) => <span className="font-medium">{s.nombre}</span> },
    {
      key: 'detalle_extra',
      header: 'Detalle',
      cell: (s) =>
        s.detalle_extra ? (
          <span className="text-muted-foreground">{s.detalle_extra}</span>
        ) : (
          <span className="text-muted-foreground">—</span>
        ),
    },
    {
      key: 'vitrina',
      header: 'Vitrina',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (s) => fmtStock(s.stock_vitrina, s.unidad),
    },
    {
      key: 'bodega',
      header: 'Bodega',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (s) => fmtStock(s.stock_bodega, s.unidad),
    },
    {
      key: 'total',
      header: 'Total',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (s) => fmtStock(s.stock_total, s.unidad),
    },
    {
      key: 'minimo',
      header: 'Mínimo',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (s) => fmtStock(s.minimo, s.unidad),
    },
    {
      key: 'estado',
      header: 'Estado',
      className: 'text-right',
      headerClassName: 'text-right',
      cell: (s) => <StockBadge total={s.stock_total} minimo={s.minimo} unidad={s.unidad} />,
    },
  ];

  if (renderActions) {
    columns.push({
      key: 'acciones',
      header: '',
      className: 'text-right',
      cell: (s) => renderActions(s),
    });
  }

  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      emptyMessage={emptyMessage}
      emptyAction={emptyAction}
      rowKey={(s) => `${s.tipo_item}-${s.item_id}`}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
