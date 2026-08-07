import { cn } from '@/lib/utils';

export function UsuarioEstadoBadge({ activo }: { activo: boolean }) {
  return (
    <span
      className={cn(
        'rounded-md px-2 py-0.5 text-xs font-medium whitespace-nowrap',
        activo ? 'bg-success/10 text-success' : 'bg-destructive/10 text-destructive',
      )}
    >
      {activo ? 'Activo' : 'Inactivo'}
    </span>
  );
}
