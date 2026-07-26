import { useMemo } from 'react';

import type { VentaItemState } from '../types';

const UMBRAL_5 = 50000;
const UMBRAL_7 = 100000;

function num(value: string | null | undefined): number {
  const n = Number.parseFloat(value ?? '');
  return Number.isNaN(n) ? 0 : n;
}

// Subtotal of one line (base by tipo_linea × cantidad, plus feromona when
// enabled). Display-only; the backend recalculates authoritatively on submit.
export function itemSubtotal(item: VentaItemState): number {
  const cantidad = item.cantidad || 0;
  let base = 0;
  switch (item.tipo_linea) {
    case 'envase_con_fragancia':
      base = num(item.envase?.precio_con_fragancia);
      break;
    case 'recarga':
      base = num(item.envase?.precio_recarga);
      break;
    case 'envase_solo':
      base = num(item.envase?.precio_solo);
      break;
    case 'producto_otro':
      base = num(item.producto?.precio);
      break;
  }
  let subtotal = base * cantidad;
  const admiteFeromona = item.tipo_linea === 'envase_con_fragancia' || item.tipo_linea === 'recarga';
  if (admiteFeromona && item.feromona_enabled && item.feromona) {
    subtotal += num(item.feromona.precio) * cantidad;
  }
  return subtotal;
}

export interface CalculoVenta {
  subtotal: number;
  descuentoPct: number;
  descuentoMonto: number;
  total: number;
}

// Mirrors backend discount_service.go: >=100000 → 7%, >=50000 → 5%, else 0.
export function useCalculoVenta(items: VentaItemState[]): CalculoVenta {
  return useMemo(() => {
    const subtotal = items.reduce((acc, item) => acc + itemSubtotal(item), 0);
    const descuentoPct = subtotal >= UMBRAL_7 ? 7 : subtotal >= UMBRAL_5 ? 5 : 0;
    const descuentoMonto = (subtotal * descuentoPct) / 100;
    const total = subtotal - descuentoMonto;
    return { subtotal, descuentoPct, descuentoMonto, total };
  }, [items]);
}
