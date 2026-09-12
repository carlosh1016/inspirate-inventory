import { format } from 'date-fns';
import { es } from 'date-fns/locale';

/** "2026-09-01" -> "Septiembre 2026". periodo es una fecha calendario (sin
 * hora), así que se parsea como fecha local para no correr de mes. */
export function formatPeriodo(periodo: string): string {
  const [year, month] = periodo.split('-').map(Number);
  const label = format(new Date(year, month - 1, 1), 'MMMM yyyy', { locale: es });
  return label.charAt(0).toUpperCase() + label.slice(1);
}
