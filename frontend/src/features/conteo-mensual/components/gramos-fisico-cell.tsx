'use client';

import { useState } from 'react';
import { Check, Loader2 } from 'lucide-react';
import { toast } from 'sonner';

import { DecimalInput } from '@/components/forms/decimal-input';
import { Button } from '@/components/ui/button';
import { formatGramos } from '@/lib/formatters';
import { getErrorMessage } from '@/lib/errors';
import { useActualizarItemConteo } from '../api/mutations';
import type { ConteoItem } from '../types';

interface Props {
  conteoId: number;
  item: ConteoItem;
  disabled: boolean;
}

// Celda editable de "gramos físico": estado local por fila + un botón
// explícito de guardar (sin auto-guardar en blur — no hay precedente de eso
// en el resto del código, todo aquí usa una acción explícita).
export function GramosFisicoCell({ conteoId, item, disabled }: Props) {
  const [value, setValue] = useState(item.gramos_fisico ?? '');
  const actualizar = useActualizarItemConteo(conteoId);
  const dirty = value.trim() !== '' && value !== (item.gramos_fisico ?? '');

  if (disabled) {
    return (
      <span className="tabular-nums">
        {item.gramos_fisico !== null ? formatGramos(item.gramos_fisico) : <span className="text-muted-foreground">—</span>}
      </span>
    );
  }

  const guardar = async () => {
    try {
      await actualizar.mutateAsync({ itemId: item.id, gramosFisico: value });
      toast.success(`${item.fragancia_nombre}: gramos físicos guardados`);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <div className="flex items-center gap-2">
      <DecimalInput value={value} onChange={setValue} suffix="g" className="w-28" />
      <Button
        size="icon"
        variant="outline"
        disabled={!dirty || actualizar.isPending}
        onClick={guardar}
        aria-label={`Guardar gramos físicos de ${item.fragancia_nombre}`}
      >
        {actualizar.isPending ? <Loader2 className="size-4 animate-spin" /> : <Check className="size-4" />}
      </Button>
    </div>
  );
}
