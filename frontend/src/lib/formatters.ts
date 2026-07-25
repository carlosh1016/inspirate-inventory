import { format, formatDistanceToNow } from 'date-fns';
import { es } from 'date-fns/locale';

/** COP without decimals and thousands separators, e.g. "$17.000". */
export function formatCurrency(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (Number.isNaN(num)) return '$0';
  return new Intl.NumberFormat('es-CO', {
    style: 'currency',
    currency: 'COP',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(num);
}

/** Grams with up to one decimal (comma), e.g. "170g" or "170,5g". */
export function formatGramos(value: string | number): string {
  const num = typeof value === 'string' ? parseFloat(value) : value;
  if (Number.isNaN(num)) return '0g';
  const rounded = Math.round(num * 10) / 10;
  return `${rounded.toString().replace('.', ',')}g`;
}

export function formatDate(iso: string): string {
  try {
    return format(new Date(iso), "d 'de' MMMM 'de' yyyy", { locale: es });
  } catch {
    return iso;
  }
}

export function formatDateShort(iso: string): string {
  try {
    return format(new Date(iso), 'd MMM yyyy', { locale: es });
  } catch {
    return iso;
  }
}

export function formatDateTime(iso: string): string {
  try {
    return format(new Date(iso), "d 'de' MMMM, h:mm a", { locale: es });
  } catch {
    return iso;
  }
}

export function formatRelative(iso: string): string {
  try {
    return formatDistanceToNow(new Date(iso), { addSuffix: true, locale: es });
  } catch {
    return iso;
  }
}

/** "V-000567" */
export function formatConsecutivoVenta(id: number): string {
  return `V-${id.toString().padStart(6, '0')}`;
}
