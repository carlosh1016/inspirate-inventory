import type { Genero, StockSnapshot } from '@/types/domain';

// Mirrors handlers/fragancias/dto.go FraganciaResponse.
export interface Fragancia {
  id: number;
  nombre_comercial: string;
  nombre_alternativo: string | null;
  genero: Genero;
  numero_genero: number;
  gramos_minimo: string;
  activo: boolean;
  stock: StockSnapshot;
  created_at: string;
  updated_at: string;
}

// A `type` (not interface) so it satisfies the Record<string, unknown>
// constraint of useUrlFilters.
export type FraganciasFilters = {
  page: number;
  q: string;
  genero: string;
  activo: string;
  stock_bajo: boolean;
};
