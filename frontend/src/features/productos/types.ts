import type { CategoriaProducto, StockSnapshot } from '@/types/domain';

export interface Producto {
  id: number;
  nombre: string;
  categoria: CategoriaProducto;
  precio: string;
  stock_minimo: number;
  activo: boolean;
  stock: StockSnapshot;
  created_at: string;
  updated_at: string;
}

export type ProductosFilters = {
  page: number;
  q: string;
  categoria: string;
  activo: string;
  stock_bajo: boolean;
};

export const CATEGORIA_OPTIONS = [
  { value: 'crema', label: 'Crema' },
  { value: 'feromona', label: 'Feromona' },
  { value: 'hogar', label: 'Hogar' },
  { value: 'regalo', label: 'Regalo' },
  { value: 'otro', label: 'Otro' },
];

export function categoriaLabel(categoria: string): string {
  return CATEGORIA_OPTIONS.find((o) => o.value === categoria)?.label ?? categoria;
}
