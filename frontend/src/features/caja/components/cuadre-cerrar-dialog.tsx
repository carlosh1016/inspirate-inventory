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
import { Textarea } from '@/components/ui/textarea';
import { formatCurrency } from '@/lib/formatters';
import { getErrorMessage } from '@/lib/errors';
import { useCerrarCuadre } from '../api/mutations';
import type { Cuadre } from '../types';

interface Props {
  cuadre: Cuadre;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CuadreCerrarDialog({ cuadre, open, onOpenChange }: Props) {
  const cerrar = useCerrarCuadre(cuadre.id);
  const [valorTurno, setValorTurno] = useState('0');
  const [observaciones, setObservaciones] = useState('');

  const submit = async () => {
    try {
      await cerrar.mutateAsync({
        valor_turno: valorTurno || '0',
        observaciones: observaciones.trim() || null,
      });
      toast.success('Caja cerrada correctamente');
      onOpenChange(false);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  const resumen: { label: string; value: string }[] = [
    { label: 'Total ventas', value: formatCurrency(cuadre.total_ventas) },
    { label: 'Total pagos', value: formatCurrency(cuadre.total_pagos) },
    { label: 'Consignaciones', value: formatCurrency(cuadre.total_consignaciones) },
    { label: 'Saldo esperado', value: formatCurrency(cuadre.saldo_calculado) },
  ];

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Cerrar caja del día</DialogTitle>
          <DialogDescription>
            Una vez cerrada, no podrás registrar más ventas ni cambios en esta caja.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <Label>Valor turno (pago a vendedora)</Label>
            <MoneyInput value={valorTurno} onChange={setValorTurno} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="cerrar-obs">Observaciones (opcional)</Label>
            <Textarea id="cerrar-obs" value={observaciones} onChange={(e) => setObservaciones(e.target.value)} maxLength={1000} />
          </div>
          <dl className="space-y-1 rounded-md bg-muted/50 p-3 text-sm">
            {resumen.map((row) => (
              <div key={row.label} className="flex justify-between">
                <dt className="text-muted-foreground">{row.label}</dt>
                <dd className="tabular-nums">{row.value}</dd>
              </div>
            ))}
          </dl>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={cerrar.isPending}>
            Cancelar
          </Button>
          <Button variant="destructive" onClick={submit} disabled={cerrar.isPending}>
            {cerrar.isPending ? 'Cerrando…' : 'Cerrar caja definitivamente'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
