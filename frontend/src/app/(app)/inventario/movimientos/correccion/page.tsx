'use client';

import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { CorreccionForm } from '@/features/movimientos/components/correccion-form';

export default function CorreccionPage() {
  return (
    <RequireRole role="admin">
      <div className="max-w-2xl">
        <PageHeader
          title="Corrección de stock"
          description="Corrige el stock a una cantidad absoluta."
          backHref="/inventario/movimientos"
          backLabel="Volver a movimientos"
        />
        <CorreccionForm />
      </div>
    </RequireRole>
  );
}
