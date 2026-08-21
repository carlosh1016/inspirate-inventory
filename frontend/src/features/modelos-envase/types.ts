export interface ModeloEnvase {
  id: number;
  tipo: string;
  tamano_oz: string;
  equiv_gramos: string;
  precio_solo: string;
  precio_con_fragancia: string;
  precio_recarga: string;
  activo: boolean;
  tiene_variantes: boolean;
  variantes_activas: number;
  created_at: string;
  updated_at: string;
}

export type ModelosEnvaseFilters = {
  page: number;
  q: string;
  activo: string;
};
