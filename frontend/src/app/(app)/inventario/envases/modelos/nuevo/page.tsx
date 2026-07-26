'use client';

import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { ModeloEnvaseForm } from '@/features/modelos-envase/components/modelo-envase-form';

export default function NuevoModeloEnvasePage() {
  return (
    <RequireRole role="admin">
      <div className="max-w-2xl">
        <PageHeader
          title="Nuevo modelo de envase"
          description="Define un tipo de envase con su tamaño y precios."
          backHref="/inventario/envases/modelos"
          backLabel="Volver a modelos"
        />
        <ModeloEnvaseForm />
      </div>
    </RequireRole>
  );
}
