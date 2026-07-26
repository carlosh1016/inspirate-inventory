'use client';

import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { formatCurrency } from '@/lib/formatters';

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  total: number;
  itemsCount: number;
  metodoPagoNombre: string;
  onConfirm: () => void;
  submitting: boolean;
}

// Guardrail for large sales (> 200.000 COP) before sending.
export function VentaConfirmarDialog({
  open,
  onOpenChange,
  total,
  itemsCount,
  metodoPagoNombre,
  onConfirm,
  submitting,
}: Props) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent showCloseButton={false}>
        <DialogHeader>
          <DialogTitle>¿Confirmas esta venta?</DialogTitle>
          <DialogDescription>Revisa el resumen antes de registrar.</DialogDescription>
        </DialogHeader>
        <dl className="space-y-1 text-sm">
          <div className="flex justify-between">
            <dt className="text-muted-foreground">Total</dt>
            <dd className="text-lg font-semibold">{formatCurrency(total)}</dd>
          </div>
          <div className="flex justify-between">
            <dt className="text-muted-foreground">Ítems</dt>
            <dd>{itemsCount}</dd>
          </div>
          <div className="flex justify-between">
            <dt className="text-muted-foreground">Método</dt>
            <dd>{metodoPagoNombre}</dd>
          </div>
        </dl>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={submitting}>
            Volver
          </Button>
          <Button onClick={onConfirm} disabled={submitting}>
            {submitting ? 'Registrando…' : 'Sí, registrar'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
