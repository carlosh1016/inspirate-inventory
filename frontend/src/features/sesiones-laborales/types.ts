export interface Sesion {
  id: number;
  usuario?: { id: number; nombre_completo: string } | null;
  entrada_at: string;
  salida_at: string | null;
  // "HH:MM:SS" or "Nd HH:MM:SS"; null while the turno is open.
  horas_trabajadas: string | null;
}

// Parses a horas_trabajadas string to total seconds.
export function horasToSeconds(horas: string): number {
  let days = 0;
  let rest = horas;
  const dayMatch = horas.match(/^(\d+)d\s+(.*)$/);
  if (dayMatch) {
    days = Number.parseInt(dayMatch[1], 10);
    rest = dayMatch[2];
  }
  const [h, m, s] = rest.split(':').map((n) => Number.parseInt(n, 10) || 0);
  return days * 86400 + h * 3600 + m * 60 + s;
}

// "6h 32m" from total seconds.
export function formatDuracion(totalSeconds: number): string {
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  return `${hours}h ${minutes.toString().padStart(2, '0')}m`;
}
