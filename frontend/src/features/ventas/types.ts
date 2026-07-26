export type TipoLinea = 'envase_con_fragancia' | 'recarga' | 'envase_solo' | 'producto_otro';

// ---- API response shapes (handlers/ventas/dto.go) ----

export interface VentaUsuarioBrief {
  id: number;
  nombre_completo: string;
}

export interface VentaMetodoPagoBrief {
  id: number;
  nombre: string;
  codigo: string;
}

export interface VentaItemDetalle {
  id: number;
  tipo_linea: TipoLinea;
  fragancia_id?: number;
  fragancia_nombre?: string;
  variante_envase_id?: number;
  variante_envase_nombre?: string;
  producto_id?: number;
  producto_nombre?: string;
  feromona_producto_id?: number;
  feromona_nombre?: string;
  gramos_fragancia?: string;
  cantidad: number;
  precio_unitario: string;
  subtotal: string;
}

export interface VentaDetallada {
  id: number;
  consecutivo: string;
  usuario: VentaUsuarioBrief;
  metodo_pago: VentaMetodoPagoBrief;
  items: VentaItemDetalle[];
  subtotal: string;
  descuento_pct: string;
  descuento_monto: string;
  total: string;
  observaciones: string | null;
  created_at: string;
  movimientos_generados: number[];
}

export interface VentaListItem {
  id: number;
  consecutivo: string;
  usuario_id: number;
  usuario_nombre: string;
  metodo_pago_id: number;
  metodo_pago_nombre: string;
  metodo_pago_codigo: string;
  items_count: number;
  subtotal: string;
  descuento_pct: string;
  descuento_monto: string;
  total: string;
  observaciones?: string | null;
  created_at: string;
}

export type VentasFilters = {
  page: number;
  metodo_pago_id: number;
  usuario_id: number;
  fecha_desde: string;
  fecha_hasta: string;
  con_descuento: boolean;
  total_min: string;
  total_max: string;
};

// ---- Local builder state (Nueva venta) ----

export interface EnvaseSeleccionado {
  variante_id: number;
  label: string;
  equiv_gramos: string;
  precio_solo: string;
  precio_con_fragancia: string;
  precio_recarga: string;
}

export interface RefSeleccionada {
  id: number;
  label: string;
}

export interface ProductoSeleccionado extends RefSeleccionada {
  precio: string;
}

export interface VentaItemState {
  key: string;
  tipo_linea: TipoLinea;
  envase: EnvaseSeleccionado | null;
  fragancia: RefSeleccionada | null;
  producto: ProductoSeleccionado | null;
  gramos: string;
  cantidad: number;
  feromona_enabled: boolean;
  feromona: ProductoSeleccionado | null;
}
