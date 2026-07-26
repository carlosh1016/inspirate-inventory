// Turns a filters object into a flat query-param map for the API: drops
// undefined/null/'' and `false` booleans (a false flag means "no filter"),
// stringifies the rest. Callers pass only the keys the endpoint accepts.
export function buildQueryParams(filters: Record<string, unknown>): Record<string, string> {
  const params: Record<string, string> = {};
  for (const [key, value] of Object.entries(filters)) {
    if (value === undefined || value === null || value === '') continue;
    if (typeof value === 'boolean') {
      if (value) params[key] = 'true';
      continue;
    }
    params[key] = String(value);
  }
  return params;
}

// The movimientos endpoint filters fecha_desde/fecha_hasta by RFC3339 instant,
// not a bare date. Convert a YYYY-MM-DD input into the start/end of that day in
// the browser's local timezone (America/Bogota for this app).
export function dateToIsoStart(date: string): string | undefined {
  if (!date) return undefined;
  const d = new Date(`${date}T00:00:00`);
  return Number.isNaN(d.getTime()) ? undefined : d.toISOString();
}

export function dateToIsoEnd(date: string): string | undefined {
  if (!date) return undefined;
  const d = new Date(`${date}T23:59:59.999`);
  return Number.isNaN(d.getTime()) ? undefined : d.toISOString();
}
