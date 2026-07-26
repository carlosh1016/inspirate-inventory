'use client';

import { RequireRole } from '@/components/guards/require-role';
import { PageHeader } from '@/components/page-header';
import { ProductoForm } from '@/features/productos/components/producto-form';

export default function NuevoProductoPage() {
  return (
    <RequireRole role="admin">
      <div className="max-w-2xl">
        <PageHeader
          title="Nuevo producto"
          description="Registra un producto en el catálogo."
          backHref="/inventario/productos"
          backLabel="Volver a productos"
        />
        <ProductoForm />
      </div>
    </RequireRole>
  );
}
