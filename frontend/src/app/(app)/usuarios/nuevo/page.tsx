'use client';

import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { UsuarioForm } from '@/features/usuarios/components/usuario-form';

export default function NuevoUsuarioPage() {
  return (
    <RequireRole role="admin">
      <PageHeader title="Nuevo usuario" backHref="/usuarios" />
      <UsuarioForm />
    </RequireRole>
  );
}
