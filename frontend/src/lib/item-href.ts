import type { TipoItem } from '@/types/domain';

// Detail page path for a catalog item, by type. Used by stock and movimientos
// tables to link a row to its item.
export function itemDetailHref(tipoItem: TipoItem, itemId: number): string {
  switch (tipoItem) {
    case 'fragancia':
      return `/inventario/fragancias/${itemId}`;
    case 'variante_envase':
      return `/inventario/envases/variantes/${itemId}`;
    case 'producto':
      return `/inventario/productos/${itemId}`;
    default:
      return '/inventario';
  }
}
