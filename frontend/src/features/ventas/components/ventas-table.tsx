'use client';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { ConsecutivoBadge } from '@/components/ventas/consecutivo-badge';
import { MetodoPagoBadge } from '@/components/ventas/metodo-pago-badge';
import { formatCurrency, formatDateTime, formatRelative } from '@/lib/formatters';
import type { Meta } from '@/types/api';
import type { VentaListItem } from '../types';

interface Props {
  data: VentaListItem[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  showVendedora: boolean;
  onPageChange: (page: number) => void;
  onRowClick: (row: VentaListItem) => void;
}

export function VentasTable({ data, meta, isLoading, page, showVendedora, onPageChange, onRowClick }: Props) {
  const columns: Column<VentaListItem>[] = [
    { key: 'consecutivo', header: 'Venta', cell: (v) => <ConsecutivoBadge value={v.consecutivo} /> },
    {
      key: 'fecha',
      header: 'Fecha',
      cell: (v) => (
        <span title={formatDateTime(v.created_at)} className="whitespace-nowrap text-muted-foreground">
          {formatRelative(v.created_at)}
        </span>
      ),
    },
    ...(showVendedora
      ? [{ key: 'vendedora', header: 'Vendedora', cell: (v: VentaListItem) => v.usuario_nombre }]
      : []),
    {
      key: 'metodo',
      header: 'Método',
      cell: (v) => <MetodoPagoBadge nombre={v.metodo_pago_nombre} codigo={v.metodo_pago_codigo} />,
    },
    { key: 'items', header: 'Ítems', className: 'text-right tabular-nums', headerClassName: 'text-right', cell: (v) => v.items_count },
    {
      key: 'subtotal',
      header: 'Subtotal',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (v) => formatCurrency(v.subtotal),
    },
    {
      key: 'descuento',
      header: 'Descuento',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (v) =>
        Number.parseFloat(v.descuento_monto) > 0 ? (
          <span className="text-success">-{formatCurrency(v.descuento_monto)}</span>
        ) : (
          <span className="text-muted-foreground">—</span>
        ),
    },
    {
      key: 'total',
      header: 'Total',
      className: 'text-right font-semibold tabular-nums',
      headerClassName: 'text-right',
      cell: (v) => formatCurrency(v.total),
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      emptyMessage="No hay ventas que coincidan con los filtros."
      rowKey={(v) => v.id}
      onRowClick={onRowClick}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
