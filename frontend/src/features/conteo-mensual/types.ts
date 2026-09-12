export type EstadoConteo = 'abierto' | 'cerrado';

export interface UsuarioBrief {
  id: number;
  nombre_completo: string;
}

// Mirrors handlers/conteos/dto.go ConteoItemResponse. Diferencia,
// excede_limite y porcentaje_vendido los calcula el backend
// (domain/conteos.ConteoItem) — el frontend nunca reimplementa la regla de
// los 2g de tolerancia.
export interface ConteoItem {
  id: number;
  fragancia_id: number;
  fragancia_nombre: string;
  saldo_inicial: string;
  gramos_sistema: string;
  gramos_fisico: string | null;
  diferencia: string | null;
  excede_limite: boolean;
  porcentaje_vendido: string | null;
}

// Mirrors handlers/conteos/dto.go ConteoResponse.
export interface Conteo {
  id: number;
  sede_id: number;
  periodo: string;
  estado: EstadoConteo;
  creado_por: UsuarioBrief;
  cerrado_por: UsuarioBrief | null;
  cerrado_at: string | null;
  created_at: string;
  items: ConteoItem[];
}

// Mirrors handlers/conteos/dto.go ConteoListItemResponse.
export interface ConteoListItem {
  id: number;
  periodo: string;
  estado: EstadoConteo;
  creado_por: UsuarioBrief;
  cerrado_por: UsuarioBrief | null;
  created_at: string;
}
