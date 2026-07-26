'use client';

import { format } from 'date-fns';
import { es } from 'date-fns/locale';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useMisSesionesMes } from '../api/use-mis-sesiones-mes';
import { formatDuracion, horasToSeconds, type Sesion } from '../types';

function hora(iso: string): string {
  try {
    return format(new Date(iso), 'h:mm a', { locale: es });
  } catch {
    return iso;
  }
}

function diaFecha(iso: string): string {
  try {
    return format(new Date(iso), "EEEE d MMM", { locale: es });
  } catch {
    return iso;
  }
}

export function HistorialMes() {
  const { data, isLoading } = useMisSesionesMes();

  const sesiones = [...(data ?? [])].sort(
    (a, b) => new Date(b.entrada_at).getTime() - new Date(a.entrada_at).getTime(),
  );

  const cerradas = sesiones.filter((s): s is Sesion & { horas_trabajadas: string } => s.horas_trabajadas != null);
  const totalSegundos = cerradas.reduce((acc, s) => acc + horasToSeconds(s.horas_trabajadas), 0);
  const diasTrabajados = new Set(sesiones.map((s) => s.entrada_at.slice(0, 10))).size;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Historial del mes</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <p className="text-sm text-muted-foreground">Cargando…</p>
        ) : sesiones.length === 0 ? (
          <p className="text-sm text-muted-foreground">Aún no tienes turnos este mes.</p>
        ) : (
          <>
            <ul className="divide-y divide-border text-sm">
              {sesiones.map((s) => (
                <li key={s.id} className="flex items-center justify-between gap-3 py-2">
                  <span className="capitalize">{diaFecha(s.entrada_at)}</span>
                  <span className="text-muted-foreground">
                    {hora(s.entrada_at)} {s.salida_at ? `- ${hora(s.salida_at)}` : ''}
                  </span>
                  <span className="tabular-nums">
                    {s.horas_trabajadas ? (
                      formatDuracion(horasToSeconds(s.horas_trabajadas))
                    ) : (
                      <span className="text-success">en curso</span>
                    )}
                  </span>
                </li>
              ))}
            </ul>
            <p className="mt-3 border-t border-border pt-3 text-sm font-medium">
              Total del mes: {formatDuracion(totalSegundos)} — {diasTrabajados}{' '}
              {diasTrabajados === 1 ? 'día trabajado' : 'días trabajados'}
            </p>
          </>
        )}
      </CardContent>
    </Card>
  );
}
