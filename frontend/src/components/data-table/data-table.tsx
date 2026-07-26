'use client';

import type { ReactNode } from 'react';

import { Skeleton } from '@/components/ui/skeleton';
import { cn } from '@/lib/utils';
import { DataTablePagination } from './data-table-pagination';
import type { Column } from './types';

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  isLoading?: boolean;
  emptyMessage?: string;
  /** Optional call-to-action rendered in the empty state (e.g. "+ Nueva"). */
  emptyAction?: ReactNode;
  totalItems?: number;
  page?: number;
  pageSize?: number;
  onPageChange?: (page: number) => void;
  rowKey: (row: T) => string | number;
  onRowClick?: (row: T) => void;
}

export function DataTable<T>({
  columns,
  data,
  isLoading = false,
  emptyMessage = 'No hay datos para mostrar.',
  emptyAction,
  totalItems,
  page,
  pageSize,
  onPageChange,
  rowKey,
  onRowClick,
}: DataTableProps<T>) {
  const showPagination =
    typeof totalItems === 'number' &&
    typeof page === 'number' &&
    typeof pageSize === 'number' &&
    !!onPageChange;

  return (
    <div>
      <div className="overflow-x-auto rounded-lg border border-border">
        <table className="w-full border-collapse text-sm">
          <thead className="border-b border-border bg-muted/40">
            <tr>
              {columns.map((column) => (
                <th
                  key={column.key}
                  className={cn(
                    'px-3 py-2.5 text-left text-xs font-medium tracking-wide text-muted-foreground uppercase',
                    column.headerClassName,
                  )}
                >
                  {column.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              Array.from({ length: 5 }).map((_, rowIndex) => (
                <tr key={`skeleton-${rowIndex}`} className="border-b border-border last:border-0">
                  {columns.map((column) => (
                    <td key={column.key} className="px-3 py-3">
                      <Skeleton className="h-4 w-full max-w-32" />
                    </td>
                  ))}
                </tr>
              ))
            ) : data.length === 0 ? (
              <tr>
                <td colSpan={columns.length} className="px-3 py-12">
                  <div className="flex flex-col items-center gap-3 text-center">
                    <p className="text-sm text-muted-foreground">{emptyMessage}</p>
                    {emptyAction}
                  </div>
                </td>
              </tr>
            ) : (
              data.map((row) => (
                <tr
                  key={rowKey(row)}
                  onClick={onRowClick ? () => onRowClick(row) : undefined}
                  className={cn(
                    'border-b border-border last:border-0',
                    onRowClick && 'cursor-pointer hover:bg-muted/40',
                  )}
                >
                  {columns.map((column) => (
                    <td key={column.key} className={cn('px-3 py-2.5 align-middle', column.className)}>
                      {column.cell(row)}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
      {showPagination && (
        <DataTablePagination
          page={page}
          pageSize={pageSize}
          totalItems={totalItems}
          onPageChange={onPageChange}
        />
      )}
    </div>
  );
}
