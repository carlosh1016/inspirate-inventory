export type EstadoCuadre = 'abierto' | 'cerrado';

export interface UsuarioBrief {
  id: number;
  nombre_completo: string;
}

export interface PagoCaja {
  id: number;
  concepto: string;
  monto: string;
  usuario: UsuarioBrief | null;
  created_at: string;
}

export interface Consignacion {
  id: number;
  monto: string;
  banco: string | null;
  referencia: string | null;
  usuario: UsuarioBrief | null;
  created_at: string;
}

export interface Cuadre {
  id: number;
  fecha: string;
  estado: EstadoCuadre;
  fondo_base: string;
  total_efectivo: string;
  total_nequi: string;
  total_daviplata: string;
  total_transferencia: string;
  total_otros: string;
  total_ventas: string;
  total_pagos: string;
  total_consignaciones: string;
  valor_turno: string;
  saldo_calculado: string;
  observaciones: string | null;
  pagos: PagoCaja[];
  consignaciones: Consignacion[];
  cerrado_por: UsuarioBrief | null;
  cerrado_at: string | null;
  created_at: string;
}

export interface CuadreListItem {
  id: number;
  fecha: string;
  estado: EstadoCuadre;
  fondo_base: string;
  total_ventas: string;
  total_pagos: string;
  total_consignaciones: string;
  saldo_calculado: string;
  cerrado_por: UsuarioBrief | null;
  created_at: string;
}

export interface Warning {
  codigo: string;
  mensaje: string;
}

export type CuadresFilters = {
  page: number;
  estado: string;
  fecha_desde: string;
  fecha_hasta: string;
};
