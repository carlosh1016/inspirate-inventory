'use client';

import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { MetodoPagoForm } from '@/features/metodos-pago/components/metodo-pago-form';

export default function NuevoMetodoPagoPage() {
  return (
    <RequireRole role="admin">
      <div className="max-w-lg">
        <PageHeader
          title="Nuevo método de pago"
          backHref="/configuracion/metodos-pago"
          backLabel="Volver a métodos de pago"
        />
        <MetodoPagoForm />
      </div>
    </RequireRole>
  );
}
