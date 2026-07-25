// App-wide constants.

/** Vendedora inactivity window before auto-logout (1 hour). */
export const INACTIVITY_TIMEOUT_MS = 60 * 60 * 1000;

/** Public routes that never require an authenticated session. */
export const PUBLIC_PATH_PREFIXES = ['/login', '/forgot-password', '/reset-password'] as const;

export function isPublicPath(pathname: string): boolean {
  return PUBLIC_PATH_PREFIXES.some((p) => pathname === p || pathname.startsWith(p + '/'));
}
