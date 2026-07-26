'use client';

import { PageHeader } from '@/components/page-header';
import { NuevaVentaForm } from '@/features/ventas/components/nueva-venta-form';

export default function NuevaVentaPage() {
  return (
    <div className="max-w-3xl">
      <PageHeader
        title="Nueva venta"
        description="Registra una venta y descuenta el inventario automáticamente."
        backHref="/ventas"
        backLabel="Volver a ventas"
      />
      <NuevaVentaForm />
    </div>
  );
}
