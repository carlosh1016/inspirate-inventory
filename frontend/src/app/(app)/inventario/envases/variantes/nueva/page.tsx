'use client';

import { PageHeader } from '@/components/page-header';
import { VarianteEnvaseForm } from '@/features/variantes-envase/components/variante-envase-form';

export default function NuevaVarianteEnvasePage() {
  return (
    <div className="max-w-2xl">
      <PageHeader
        title="Nueva variante de envase"
        description="Registra un color para un modelo de envase."
        backHref="/inventario/envases/variantes"
        backLabel="Volver a variantes"
      />
      <VarianteEnvaseForm />
    </div>
  );
}
