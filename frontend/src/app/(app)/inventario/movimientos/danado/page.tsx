'use client';

import { PageHeader } from '@/components/page-header';
import { DanadoForm } from '@/features/movimientos/components/danado-form';

export default function DanadoPage() {
  return (
    <div className="max-w-2xl">
      <PageHeader
        title="Producto dañado"
        description="Registra una baja por producto dañado."
        backHref="/inventario/movimientos"
        backLabel="Volver a movimientos"
      />
      <DanadoForm />
    </div>
  );
}
