'use client';

import { useState } from 'react';
import { toast } from 'sonner';

import { MoneyInput } from '@/components/forms/money-input';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { getErrorMessage } from '@/lib/errors';
import { useAbrirCuadre } from '../api/mutations';

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CuadreAbrirDialog({ open, onOpenChange }: Props) {
  const abrir = useAbrirCuadre();
  const [fondoBase, setFondoBase] = useState('100000');

  const submit = async () => {
    try {
      const result = await abrir.mutateAsync(fondoBase || null);
      toast.success('Caja abierta');
      result.warnings.forEach((w) => toast.warning(w.mensaje));
      onOpenChange(false);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Abrir caja del día</DialogTitle>
          <DialogDescription>Define el fondo base con el que inicia la caja.</DialogDescription>
        </DialogHeader>
        <div className="space-y-2">
          <Label>Fondo base</Label>
          <MoneyInput value={fondoBase} onChange={setFondoBase} />
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={abrir.isPending}>
            Cancelar
          </Button>
          <Button onClick={submit} disabled={abrir.isPending}>
            {abrir.isPending ? 'Abriendo…' : 'Abrir caja'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
