'use client';

import Link from 'next/link';
import { Pencil, Power } from 'lucide-react';

import { DataTable } from '@/components/data-table/data-table';
import type { Column } from '@/components/data-table/types';
import { Button } from '@/components/ui/button';
import { formatDate } from '@/lib/formatters';
import type { Meta } from '@/types/api';
import type { UsuarioApi } from '@/types/domain';
import { UsuarioEstadoBadge } from './usuario-estado-badge';
import { UsuarioRolBadge } from './usuario-rol-badge';

interface Props {
  data: UsuarioApi[];
  meta?: Meta;
  isLoading: boolean;
  page: number;
  onPageChange: (page: number) => void;
  currentUserId: number;
  onToggleEstado: (usuario: UsuarioApi) => void;
}

export function UsuariosTable({ data, meta, isLoading, page, onPageChange, currentUserId, onToggleEstado }: Props) {
  const columns: Column<UsuarioApi>[] = [
    { key: 'nombre', header: 'Nombre', cell: (u) => <span className="font-medium">{u.nombre_completo}</span> },
    { key: 'correo', header: 'Correo', cell: (u) => <span className="text-muted-foreground">{u.correo}</span> },
    { key: 'rol', header: 'Rol', cell: (u) => <UsuarioRolBadge rol={u.rol} /> },
    { key: 'estado', header: 'Estado', cell: (u) => <UsuarioEstadoBadge activo={u.is_active} /> },
    { key: 'creado', header: 'Creado', cell: (u) => formatDate(u.created_at) },
    {
      key: 'acciones',
      header: '',
      className: 'text-right',
      cell: (u) =>
        u.id === currentUserId ? null : (
          <div className="flex justify-end gap-1">
            <Link href={`/usuarios/${u.id}/editar`}>
              <Button variant="ghost" size="icon-sm" aria-label="Editar usuario">
                <Pencil className="size-4" />
              </Button>
            </Link>
            <Button
              variant="ghost"
              size="icon-sm"
              aria-label={u.is_active ? 'Desactivar usuario' : 'Activar usuario'}
              onClick={() => onToggleEstado(u)}
            >
              <Power className={u.is_active ? 'size-4 text-destructive' : 'size-4 text-success'} />
            </Button>
          </div>
        ),
    },
  ];

  return (
    <DataTable
      columns={columns}
      data={data}
      isLoading={isLoading}
      emptyMessage="No hay usuarios que coincidan con los filtros."
      rowKey={(u) => u.id}
      totalItems={meta?.total}
      page={page}
      pageSize={meta?.page_size}
      onPageChange={onPageChange}
    />
  );
}
