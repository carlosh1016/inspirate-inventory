import type { StockSnapshot } from '@/types/domain';

export interface VarianteEnvase {
  id: number;
  modelo_envase_id: number;
  color: string;
  stock_minimo: number;
  activo: boolean;
  stock: StockSnapshot;
  created_at: string;
  updated_at: string;
}

export type VariantesEnvaseFilters = {
  page: number;
  q: string;
  modelo_envase_id: number;
  activo: string;
  stock_bajo: boolean;
};
