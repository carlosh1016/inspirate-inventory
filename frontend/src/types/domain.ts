// Domain types shared across features. Shapes mirror the backend DTOs
// (verified against the Go handlers), using snake_case for JSON fidelity.

export type Rol = 'admin' | 'vendedora';

export type Genero = 'masculina' | 'femenina';
export type CategoriaProducto = 'crema' | 'feromona' | 'hogar' | 'regalo' | 'otro';
export type TipoItem = 'fragancia' | 'variante_envase' | 'producto';
export type Ubicacion = 'vitrina' | 'bodega';

/** vitrina/bodega/total stock snapshot (decimal strings) shared by catalog items. */
export interface StockSnapshot {
  vitrina: string;
  bodega: string;
  total: string;
}

/** The authenticated user kept in the auth store. */
export interface UsuarioSession {
  id: number;
  nombre_completo: string;
  correo: string;
  rol: Rol;
  sede_id: number;
}

/** Full usuario shape returned by /auth/login and /auth/me. */
export interface UsuarioApi extends UsuarioSession {
  is_active: boolean;
  last_login_at: string | null;
  created_at: string;
}

export function toUsuarioSession(u: UsuarioApi): UsuarioSession {
  return {
    id: u.id,
    nombre_completo: u.nombre_completo,
    correo: u.correo,
    rol: u.rol,
    sede_id: u.sede_id,
  };
}
