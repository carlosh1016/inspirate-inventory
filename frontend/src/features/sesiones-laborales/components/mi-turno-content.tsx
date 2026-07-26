'use client';

import Link from 'next/link';

import { ErrorState } from '@/components/feedback/error-state';
import { LoadingState } from '@/components/feedback/loading-state';
import { PageHeader } from '@/components/page-header';
import { buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { usePermission } from '@/hooks/use-permission';
import { useAuthStore } from '@/stores/auth-store';
import { useMiSesionAbierta } from '../api/use-mi-sesion-abierta';
import { HistorialMes } from './historial-mes';
import { TurnoAbiertoCard } from './turno-abierto-card';
import { TurnoCerradoCard } from './turno-cerrado-card';

export function MiTurnoContent() {
  const { isAdmin } = usePermission();
  const usuario = useAuthStore((s) => s.usuario);
  const { data: sesion, isLoading, isError, error, refetch } = useMiSesionAbierta();

  if (isAdmin) {
    return (
      <>
        <PageHeader title="Mi turno" />
        <Card>
          <CardContent className="flex flex-col items-center gap-4 py-16 text-center">
            <p className="text-sm text-muted-foreground">Esta pantalla es solo para vendedoras.</p>
            <Link href="/dashboard" className={buttonVariants({ variant: 'outline' })}>
              Volver al dashboard
            </Link>
          </CardContent>
        </Card>
      </>
    );
  }

  return (
    <>
      <PageHeader title={`Mi turno${usuario ? ` · ${usuario.nombre_completo}` : ''}`} />
      {isLoading ? (
        <LoadingState />
      ) : isError ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : (
        <div className="space-y-6">
          {sesion ? <TurnoAbiertoCard sesion={sesion} /> : <TurnoCerradoCard />}
          <HistorialMes />
        </div>
      )}
    </>
  );
}
