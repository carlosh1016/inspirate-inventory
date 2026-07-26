'use client';

import { PageHeader } from '@/components/page-header';
import { EntradaMercanciaForm } from '@/features/movimientos/components/entrada-mercancia-form';

export default function EntradaMercanciaPage() {
  return (
    <div className="max-w-3xl">
      <PageHeader
        title="Entrada de mercancía"
        description="Registra la mercancía que llega del proveedor."
        backHref="/inventario/movimientos"
        backLabel="Volver a movimientos"
      />
      <EntradaMercanciaForm />
    </div>
  );
}
