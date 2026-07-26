'use client';

import Link from 'next/link';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { formatConsecutivoVenta, formatDateTime, formatGramos, formatRelative } from '@/lib/formatters';
import { itemDetailHref } from '@/lib/item-href';
import { cn } from '@/lib/utils';
import type { Meta } from '@/types/api';
import { TIPO_MOVIMIENTO_META } from '../tipo-meta';
import type { Movimiento } from '../types';

interface Props {
  data: Movimiento[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  onPageChange: (page: number) => void;
}

function fmtCantidad(mov: Movimiento): string {
  const n = Number.parseFloat(mov.cantidad);
  const sign = n > 0 ? '+' : '';
  const body = mov.tipo_item === 'fragancia' ? formatGramos(n) : String(n);
  return `${sign}${body}`;
}

const columns: Column<Movimiento>[] = [
  {
    key: 'fecha',
    header: 'Fecha',
    cell: (m) => (
      <span title={formatDateTime(m.created_at)} className="whitespace-nowrap text-muted-foreground">
        {formatRelative(m.created_at)}
      </span>
    ),
  },
  {
    key: 'tipo',
    header: 'Tipo',
    cell: (m) => {
      const meta = TIPO_MOVIMIENTO_META[m.tipo];
      return (
        <span className={cn('rounded-md px-2 py-0.5 text-xs font-medium whitespace-nowrap', meta.className)}>
          {meta.label}
        </span>
      );
    },
  },
  {
    key: 'item',
    header: 'Ítem',
    cell: (m) => (
      <Link
        href={itemDetailHref(m.tipo_item, m.item.id)}
        className="font-medium hover:underline"
        onClick={(e) => e.stopPropagation()}
      >
        {m.item.nombre}
      </Link>
    ),
  },
  {
    key: 'cantidad',
    header: 'Cantidad',
    className: 'text-right tabular-nums',
    headerClassName: 'text-right',
    cell: (m) => {
      const negative = Number.parseFloat(m.cantidad) < 0;
      return <span className={negative ? 'text-destructive' : 'text-success'}>{fmtCantidad(m)}</span>;
    },
  },
  {
    key: 'stock',
    header: 'Anterior → Posterior',
    className: 'text-right tabular-nums whitespace-nowrap',
    headerClassName: 'text-right',
    cell: (m) => `${m.stock_anterior} → ${m.stock_posterior}`,
  },
  { key: 'ubicacion', header: 'Ubicación', cell: (m) => <span className="capitalize">{m.ubicacion}</span> },
  {
    key: 'usuario',
    header: 'Usuario',
    cell: (m) => m.usuario.nombre_completo ?? <span className="text-muted-foreground">—</span>,
  },
  {
    key: 'motivo',
    header: 'Motivo',
    cell: (m) =>
      m.motivo ? (
        <span title={m.motivo} className="block max-w-40 truncate text-muted-foreground">
          {m.motivo}
        </span>
      ) : (
        <span className="text-muted-foreground">—</span>
      ),
  },
  {
    key: 'venta',
    header: 'Venta',
    cell: (m) =>
      m.venta_id ? (
        <Link
          href={`/ventas/${m.venta_id}`}
          className="hover:underline"
          onClick={(e) => e.stopPropagation()}
        >
          {formatConsecutivoVenta(m.venta_id)}
        </Link>
      ) : (
        <span className="text-muted-foreground">—</span>
      ),
  },
];

export function MovimientosTable({ data, meta, isLoading, page, onPageChange }: Props) {
  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      emptyMessage="No hay movimientos registrados en el período."
      rowKey={(m) => m.id}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
