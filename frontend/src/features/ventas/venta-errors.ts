import { parseApiError } from '@/lib/errors';

// POST /ventas 422s come in three shapes (see domain/ventas/errors.go):
//  - stock insuficiente:  extra.items[] = {index, ..., requerido, disponible, unidad}
//  - coherencia de item:  extra.items[] = {index, motivo}
//  - cuadre cerrado:       business_rule, title "Cuadre de caja cerrado", no extra
// All item errors are keyed by the request array `index`.

export type VentaErrorKind = 'stock' | 'coherence' | 'cuadre_cerrado' | 'other';

export interface VentaErrorInfo {
  kind: VentaErrorKind;
  message: string;
  /** index -> message to show under the affected item row. */
  byIndex: Map<number, string>;
}

export function parseVentaError(error: unknown): VentaErrorInfo {
  const problem = parseApiError(error);
  const extra = problem.extra as { items?: unknown } | undefined;
  const items = extra && Array.isArray(extra.items) ? (extra.items as Record<string, unknown>[]) : null;

  if (items && items.length > 0) {
    const byIndex = new Map<number, string>();
    const isStock = 'disponible' in items[0];
    for (const it of items) {
      const index = Number(it.index ?? -1);
      if (isStock) {
        const disponible = String(it.disponible ?? '');
        const suffix = it.unidad === 'gramos' ? 'g' : '';
        byIndex.set(index, `Stock insuficiente. Disponible: ${disponible}${suffix}`);
      } else {
        byIndex.set(index, String(it.motivo ?? 'Revisa este ítem'));
      }
    }
    return {
      kind: isStock ? 'stock' : 'coherence',
      message: isStock
        ? 'Uno o más ítems no tienen stock suficiente.'
        : (problem.detail ?? 'Revisa los ítems marcados.'),
      byIndex,
    };
  }

  if (problem.title === 'Cuadre de caja cerrado') {
    return { kind: 'cuadre_cerrado', message: problem.detail ?? problem.title, byIndex: new Map() };
  }

  return { kind: 'other', message: problem.detail ?? problem.title, byIndex: new Map() };
}
