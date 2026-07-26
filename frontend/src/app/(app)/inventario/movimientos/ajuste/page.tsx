'use client';

import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { AjusteForm } from '@/features/movimientos/components/ajuste-form';

export default function AjustePage() {
  return (
    <RequireRole role="admin">
      <div className="max-w-2xl">
        <PageHeader
          title="Ajuste de stock"
          description="Corrige el stock a una cantidad absoluta."
          backHref="/inventario/movimientos"
          backLabel="Volver a movimientos"
        />
        <AjusteForm />
      </div>
    </RequireRole>
  );
}
