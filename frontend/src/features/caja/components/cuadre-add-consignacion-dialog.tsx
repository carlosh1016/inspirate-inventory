'use client';

import { useState } from 'react';
import { toast } from 'sonner';

import { MoneyInput } from '@/components/forms/money-input';
import { FormField } from '@/components/forms/form-field';
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
import { useAddConsignacion } from '../api/mutations';

interface Props {
  cuadreId: number;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CuadreAddConsignacionDialog({ cuadreId, open, onOpenChange }: Props) {
  const addConsignacion = useAddConsignacion(cuadreId);
  const [monto, setMonto] = useState('');
  const [banco, setBanco] = useState('');
  const [referencia, setReferencia] = useState('');
  const [error, setError] = useState<string | undefined>();

  const submit = async () => {
    if (!monto || Number.parseInt(monto, 10) <= 0) return setError('Ingresa un monto válido');
    setError(undefined);
    try {
      await addConsignacion.mutateAsync({
        monto,
        banco: banco.trim() || null,
        referencia: referencia.trim() || null,
      });
      toast.success('Consignación registrada');
      setMonto('');
      setBanco('');
      setReferencia('');
      onOpenChange(false);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Registrar consignación</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Monto</Label>
            <MoneyInput value={monto} onChange={setMonto} />
            {error && <p className="text-sm text-destructive">{error}</p>}
          </div>
          <FormField id="banco" label="Banco (opcional)">
            <Input id="banco" value={banco} onChange={(e) => setBanco(e.target.value)} />
          </FormField>
          <FormField id="referencia" label="Referencia (opcional)">
            <Input id="referencia" value={referencia} onChange={(e) => setReferencia(e.target.value)} />
          </FormField>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={addConsignacion.isPending}>
            Cancelar
          </Button>
          <Button onClick={submit} disabled={addConsignacion.isPending}>
            {addConsignacion.isPending ? 'Guardando…' : 'Registrar'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
