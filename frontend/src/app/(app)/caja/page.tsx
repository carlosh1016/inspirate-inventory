'use client';

import Link from 'next/link';

import { PageHeader } from '@/components/page-header';
import { buttonVariants } from '@/components/ui/button';
import { usePermission } from '@/hooks/use-permission';
import { CuadreHoyView } from '@/features/caja/components/cuadre-hoy-view';

export default function CajaPage() {
  const { isAdmin } = usePermission();
  return (
    <>
      <PageHeader
        title="Caja"
        description="Cuadre de caja del día."
        action={
          isAdmin ? (
            <Link href="/caja/historial" className={buttonVariants({ variant: 'outline' })}>
              Ver historial
            </Link>
          ) : undefined
        }
      />
      <CuadreHoyView />
    </>
  );
}
