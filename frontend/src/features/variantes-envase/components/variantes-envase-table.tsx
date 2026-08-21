'use client';

import type { ReactNode } from 'react';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { StockBadge } from '@/components/stock-badge';
import type { Meta } from '@/types/api';
import type { VarianteEnvase } from '../types';

interface Props {
  data: VarianteEnvase[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  modelosMap?: Map<number, string>;
  onPageChange: (page: number) => void;
  onRowClick: (row: VarianteEnvase) => void;
  emptyAction?: ReactNode;
}

export function VariantesEnvaseTable({
  data,
  meta,
  isLoading,
  page,
  modelosMap,
  onPageChange,
  onRowClick,
  emptyAction,
}: Props) {
  const columns: Column<VarianteEnvase>[] = [
    {
      key: 'modelo',
      header: 'Modelo',
      cell: (v) => (
        <span className="font-medium">
          {modelosMap?.get(v.modelo_envase_id) ?? `Modelo #${v.modelo_envase_id}`}
          {!v.activo && <span className="ml-2 text-xs text-muted-foreground">(inactiva)</span>}
        </span>
      ),
    },
    { key: 'color', header: 'Grosor', cell: (v) => v.color },
    {
      key: 'vitrina',
      header: 'Vitrina',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (v) => v.stock.vitrina,
    },
    {
      key: 'bodega',
      header: 'Bodega',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (v) => v.stock.bodega,
    },
    {
      key: 'total',
      header: 'Total',
      className: 'text-right',
      headerClassName: 'text-right',
      cell: (v) => <StockBadge total={v.stock.total} minimo={v.stock_minimo} unidad="unidades" />,
    },
    {
      key: 'stock_minimo',
      header: 'Mínimo',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (v) => v.stock_minimo,
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      emptyMessage="No hay variantes que coincidan con los filtros."
      emptyAction={emptyAction}
      rowKey={(v) => v.id}
      onRowClick={onRowClick}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
