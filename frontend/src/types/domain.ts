// Domain types shared across features. Shapes mirror the backend DTOs
// (verified against the Go handlers), using snake_case for JSON fidelity.

export type Rol = 'admin' | 'vendedora';

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
