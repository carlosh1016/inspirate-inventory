'use client';

import { useState } from 'react';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useElapsedTime } from '@/hooks/use-elapsed-time';
import { getErrorMessage } from '@/lib/errors';
import { formatDateTime, formatRelative } from '@/lib/formatters';
import { useMarcarSalida } from '../api/use-marcar-salida';
import { formatDuracion, type Sesion } from '../types';

const pad = (n: number) => n.toString().padStart(2, '0');

export function TurnoAbiertoCard({ sesion }: { sesion: Sesion }) {
  const marcarSalida = useMarcarSalida();
  const elapsed = useElapsedTime(sesion.entrada_at);
  const [confirmOpen, setConfirmOpen] = useState(false);

  const salir = async () => {
    try {
      await marcarSalida.mutateAsync();
      toast.success('Salida registrada');
      setConfirmOpen(false);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  const trabajado = elapsed ? formatDuracion(Math.floor(elapsed.totalMs / 1000)) : '';

  return (
    <Card>
      <CardContent className="flex flex-col items-center gap-4 py-12 text-center">
        <div className="space-y-1">
          <p className="text-sm text-muted-foreground">Turno actual</p>
          <p className="text-sm">
            Entrada: {formatDateTime(sesion.entrada_at)}{' '}
            <span className="text-muted-foreground">({formatRelative(sesion.entrada_at)})</span>
          </p>
        </div>
        <div>
          <p className="font-mono text-5xl font-semibold tabular-nums">
            {elapsed ? `${pad(elapsed.hours)}:${pad(elapsed.minutes)}:${pad(elapsed.seconds)}` : '00:00:00'}
          </p>
          <p className="mt-1 text-xs text-muted-foreground">tiempo trabajado</p>
        </div>
        <Button size="lg" variant="outline" onClick={() => setConfirmOpen(true)}>
          Marcar salida
        </Button>
      </CardContent>

      <Dialog open={confirmOpen} onOpenChange={setConfirmOpen}>
        <DialogContent showCloseButton={false}>
          <DialogHeader>
            <DialogTitle>¿Confirmas tu salida?</DialogTitle>
            <DialogDescription>
              Estás a punto de cerrar tu turno. Trabajaste {trabajado}.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setConfirmOpen(false)} disabled={marcarSalida.isPending}>
              Cancelar
            </Button>
            <Button onClick={salir} disabled={marcarSalida.isPending}>
              {marcarSalida.isPending ? 'Registrando…' : 'Sí, marcar salida'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  );
}
