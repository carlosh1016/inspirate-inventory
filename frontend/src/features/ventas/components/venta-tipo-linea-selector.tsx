'use client';

import { Package, RefreshCw, ShoppingBag, Sparkles, type LucideIcon } from 'lucide-react';

import { Button } from '@/components/ui/button';
import type { TipoLinea } from '../types';

const OPTIONS: { tipo: TipoLinea; label: string; icon: LucideIcon }[] = [
  { tipo: 'envase_con_fragancia', label: 'Envase con fragancia', icon: Package },
  { tipo: 'recarga', label: 'Recarga', icon: RefreshCw },
  { tipo: 'envase_solo', label: 'Envase solo', icon: ShoppingBag },
  { tipo: 'producto_otro', label: 'Otro', icon: Sparkles },
];

export function VentaTipoLineaSelector({ onAdd }: { onAdd: (tipo: TipoLinea) => void }) {
  return (
    <div className="flex flex-wrap gap-2">
      {OPTIONS.map((option) => (
        <Button key={option.tipo} type="button" variant="outline" onClick={() => onAdd(option.tipo)}>
          <option.icon className="size-4" />
          {option.label}
        </Button>
      ))}
    </div>
  );
}
