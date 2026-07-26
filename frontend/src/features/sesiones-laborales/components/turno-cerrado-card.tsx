'use client';

import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { getErrorMessage } from '@/lib/errors';
import { useMarcarEntrada } from '../api/use-marcar-entrada';

export function TurnoCerradoCard() {
  const marcarEntrada = useMarcarEntrada();

  const entrar = async () => {
    try {
      await marcarEntrada.mutateAsync();
      toast.success('Entrada registrada');
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Card>
      <CardContent className="flex flex-col items-center gap-6 py-16 text-center">
        <p className="text-lg font-medium">No tienes un turno abierto</p>
        <Button size="lg" onClick={entrar} disabled={marcarEntrada.isPending}>
          {marcarEntrada.isPending ? 'Registrando…' : 'Marcar entrada'}
        </Button>
      </CardContent>
    </Card>
  );
}
