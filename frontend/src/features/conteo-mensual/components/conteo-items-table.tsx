'use client';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { formatGramos } from '@/lib/formatters';
import { cn } from '@/lib/utils';
import type { ConteoItem } from '../types';
import { GramosFisicoCell } from './gramos-fisico-cell';

interface Props {
  conteoId: number;
  items: ConteoItem[];
  disabled: boolean;
}

export function ConteoItemsTable({ conteoId, items, disabled }: Props) {
  const columns: Column<ConteoItem>[] = [
    {
      key: 'fragancia',
      header: 'Fragancia',
      cell: (i) => <span className="font-medium">{i.fragancia_nombre}</span>,
    },
    {
      key: 'saldo_inicial',
      header: 'Saldo inicial',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (i) => formatGramos(i.saldo_inicial),
    },
    {
      key: 'gramos_sistema',
      header: 'Gramos sistema',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (i) => formatGramos(i.gramos_sistema),
    },
    {
      key: 'gramos_fisico',
      header: 'Gramos físico',
      cell: (i) => <GramosFisicoCell conteoId={conteoId} item={i} disabled={disabled} />,
    },
    {
      key: 'diferencia',
      header: 'Diferencia',
      className: 'text-right',
      headerClassName: 'text-right',
      cell: (i) =>
        i.diferencia === null ? (
          <span className="text-muted-foreground">—</span>
        ) : (
          <span className={cn('tabular-nums', i.excede_limite && 'font-semibold text-destructive')}>
            {formatGramos(i.diferencia)}
            {i.excede_limite && ' · revisar'}
          </span>
        ),
    },
    {
      key: 'porcentaje_vendido',
      header: '% vendido',
      className: 'text-right tabular-nums',
      headerClassName: 'text-right',
      cell: (i) => (i.porcentaje_vendido === null ? '—' : `${i.porcentaje_vendido}%`),
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={items}
      isLoading={false}
      emptyMessage="Este conteo no tiene fragancias activas para contar."
      rowKey={(i) => i.id}
    />
  );
}
