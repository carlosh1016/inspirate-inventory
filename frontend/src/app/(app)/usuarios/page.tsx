'use client';

import { Suspense, useState } from 'react';
import Link from 'next/link';
import { Plus } from 'lucide-react';
import { toast } from 'sonner';

import { ConfirmDeleteDialog } from '@/components/confirm-delete-dialog';
import { ErrorState } from '@/components/feedback/error-state';
import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { buttonVariants } from '@/components/ui/button';
import { useActivarUsuario } from '@/features/usuarios/api/use-activar-usuario';
import { useDesactivarUsuario } from '@/features/usuarios/api/use-desactivar-usuario';
import { useUsuarios } from '@/features/usuarios/api/use-usuarios';
import { UsuariosFiltersBar } from '@/features/usuarios/components/usuarios-filters';
import { UsuariosTable } from '@/features/usuarios/components/usuarios-table';
import type { UsuariosFilters } from '@/features/usuarios/types';
import { useUrlFilters } from '@/hooks/use-url-filters';
import { getErrorMessage } from '@/lib/errors';
import { useAuthStore } from '@/stores/auth-store';
import type { UsuarioApi } from '@/types/domain';

const NuevoButton = (
  <Link href="/usuarios/nuevo" className={buttonVariants()}>
    <Plus className="size-4" />
    Nuevo usuario
  </Link>
);

function UsuariosView() {
  const usuarioActual = useAuthStore((s) => s.usuario);
  const [toggleTarget, setToggleTarget] = useState<UsuarioApi | null>(null);
  const activar = useActivarUsuario();
  const desactivar = useDesactivarUsuario();

  const { filters, setFilter } = useUrlFilters<UsuariosFilters>({
    defaults: { page: 1, rol: 'all', activo: 'all' },
    parsers: {
      page: (v) => Math.max(1, Number(v) || 1),
      rol: (v) => (v ?? 'all') as UsuariosFilters['rol'],
      activo: (v) => (v ?? 'all') as UsuariosFilters['activo'],
    },
    serializers: {
      page: (v) => (v === 1 ? null : String(v)),
      rol: (v) => (v === 'all' ? null : v),
      activo: (v) => (v === 'all' ? null : v),
    },
  });

  const { data, isLoading, isError, error, refetch } = useUsuarios(filters);

  return (
    <>
      <PageHeader title="Usuarios" description="Cuentas de acceso al sistema." action={NuevoButton} />
      <UsuariosFiltersBar filters={filters} setFilter={setFilter} />
      {isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <UsuariosTable
          data={data?.items ?? []}
          meta={data?.meta}
          isLoading={isLoading}
          page={filters.page}
          onPageChange={(p) => setFilter('page', p)}
          currentUserId={usuarioActual?.id ?? 0}
          onToggleEstado={setToggleTarget}
        />
      )}
      <ConfirmDeleteDialog
        open={toggleTarget !== null}
        onOpenChange={(open) => !open && setToggleTarget(null)}
        title={toggleTarget?.is_active ? 'Desactivar usuario' : 'Activar usuario'}
        description={
          toggleTarget?.is_active
            ? `¿Seguro que deseas desactivar a ${toggleTarget?.nombre_completo}? No podrá iniciar sesión.`
            : `¿Seguro que deseas activar a ${toggleTarget?.nombre_completo}?`
        }
        confirmLabel={toggleTarget?.is_active ? 'Desactivar' : 'Activar'}
        onConfirm={async () => {
          if (!toggleTarget) return;
          try {
            if (toggleTarget.is_active) {
              await desactivar.mutateAsync(toggleTarget.id);
              toast.success('Usuario desactivado');
            } else {
              await activar.mutateAsync(toggleTarget.id);
              toast.success('Usuario activado');
            }
            setToggleTarget(null);
          } catch (err) {
            toast.error(getErrorMessage(err));
          }
        }}
      />
    </>
  );
}

export default function UsuariosPage() {
  return (
    <RequireRole role="admin">
      <Suspense>
        <UsuariosView />
      </Suspense>
    </RequireRole>
  );
}
