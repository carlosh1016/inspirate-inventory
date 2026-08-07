'use client';

import { use } from 'react';

import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { useUsuario } from '@/features/usuarios/api/use-usuario';
import { UsuarioForm } from '@/features/usuarios/components/usuario-form';

export default function EditarUsuarioPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const usuarioId = Number(id);
  const { data: usuario, isLoading, isError, error, refetch } = useUsuario(usuarioId);

  return (
    <RequireRole role="admin">
      <PageHeader title="Editar usuario" backHref="/usuarios" />
      {isLoading ? (
        <LoadingState />
      ) : isError || !usuario ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <UsuarioForm initialData={usuario} />
      )}
    </RequireRole>
  );
}
