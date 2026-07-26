'use client';

import { useState } from 'react';
import { Pencil } from 'lucide-react';
import { toast } from 'sonner';

import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { getErrorMessage } from '@/lib/errors';
import { usePermission } from '@/hooks/use-permission';
import { useUpdateVenta } from '../api/use-update-venta';

interface Props {
  ventaId: number;
  observaciones: string | null;
}

export function VentaObservacionesInline({ ventaId, observaciones }: Props) {
  const { isAdmin } = usePermission();
  const update = useUpdateVenta(ventaId);
  const [editing, setEditing] = useState(false);
  const [value, setValue] = useState(observaciones ?? '');

  const save = async () => {
    try {
      await update.mutateAsync(value.trim() ? value.trim() : null);
      toast.success('Observaciones actualizadas');
      setEditing(false);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  if (editing) {
    return (
      <div className="space-y-2">
        <Textarea value={value} onChange={(e) => setValue(e.target.value)} maxLength={1000} autoFocus />
        <div className="flex justify-end gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              setValue(observaciones ?? '');
              setEditing(false);
            }}
            disabled={update.isPending}
          >
            Cancelar
          </Button>
          <Button size="sm" onClick={save} disabled={update.isPending}>
            {update.isPending ? 'Guardando…' : 'Guardar'}
          </Button>
        </div>
      </div>
    );
  }

  if (!isAdmin) {
    return (
      <p className="text-sm text-muted-foreground">
        {observaciones || 'Sin observaciones.'}
      </p>
    );
  }

  return (
    <button
      type="button"
      onClick={() => setEditing(true)}
      className="group flex w-full items-start justify-between gap-2 rounded-md text-left text-sm hover:bg-muted/50"
    >
      <span className={observaciones ? '' : 'text-muted-foreground'}>
        {observaciones || 'Agregar observaciones…'}
      </span>
      <Pencil className="mt-0.5 size-3.5 shrink-0 text-muted-foreground opacity-0 group-hover:opacity-100" />
    </button>
  );
}
