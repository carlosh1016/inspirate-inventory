'use client';

import { PageHeader } from '@/components/page-header';
import { FraganciaForm } from '@/features/fragancias/components/fragancia-form';

export default function NuevaFraganciaPage() {
  return (
    <div className="max-w-2xl">
      <PageHeader
        title="Nueva fragancia"
        description="Registra una fragancia en el catálogo."
        backHref="/inventario/fragancias"
        backLabel="Volver a fragancias"
      />
      <FraganciaForm />
    </div>
  );
}
