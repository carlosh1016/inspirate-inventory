import { cn } from '@/lib/utils';

// Renders a venta consecutivo ("V-000123") as a monospace badge. Accepts the
// preformatted string or a numeric id.
export function ConsecutivoBadge({
  value,
  className,
}: {
  value: string | number;
  className?: string;
}) {
  const label = typeof value === 'number' ? `V-${value.toString().padStart(6, '0')}` : value;
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-md bg-muted px-2 py-0.5 font-mono text-xs font-medium',
        className,
      )}
    >
      {label}
    </span>
  );
}
