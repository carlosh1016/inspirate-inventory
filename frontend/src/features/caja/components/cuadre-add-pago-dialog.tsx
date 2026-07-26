'use client';

import { useState } from 'react';
import { toast } from 'sonner';

import { FormField } from '@/components/forms/form-field';
import { MoneyInput } from '@/components/forms/money-input';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { getErrorMessage } from '@/lib/errors';
import { useAddPago } from '../api/mutations';

interface Props {
  cuadreId: number;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CuadreAddPagoDialog({ cuadreId, open, onOpenChange }: Props) {
  const addPago = useAddPago(cuadreId);
  const [concepto, setConcepto] = useState('');
  const [monto, setMonto] = useState('');
  const [error, setError] = useState<string | undefined>();

  const submit = async () => {
    if (concepto.trim().length < 2) return setError('Ingresa un concepto');
    if (!monto || Number.parseInt(monto, 10) <= 0) return setError('Ingresa un monto válido');
    setError(undefined);
    try {
      await addPago.mutateAsync({ concepto: concepto.trim(), monto });
      toast.success('Pago registrado');
      setConcepto('');
      setMonto('');
      onOpenChange(false);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Registrar pago</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <FormField id="concepto" label="Concepto">
            <Input id="concepto" value={concepto} onChange={(e) => setConcepto(e.target.value)} placeholder="Ej. Papel higiénico" />
          </FormField>
          <div className="space-y-2">
            <Label>Monto</Label>
            <MoneyInput value={monto} onChange={setMonto} />
            {error && <p className="text-sm text-destructive">{error}</p>}
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={addPago.isPending}>
            Cancelar
          </Button>
          <Button onClick={submit} disabled={addPago.isPending}>
            {addPago.isPending ? 'Guardando…' : 'Registrar'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
