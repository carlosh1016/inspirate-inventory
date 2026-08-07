import { cn } from '@/lib/utils';
import type { Rol } from '@/types/domain';

const ROL_META: Record<Rol, { label: string; className: string }> = {
  admin: { label: 'Admin', className: 'bg-blue-500/10 text-blue-600 dark:text-blue-400' },
  vendedora: { label: 'Vendedora', className: 'bg-muted text-muted-foreground' },
};

export function UsuarioRolBadge({ rol }: { rol: Rol }) {
  const meta = ROL_META[rol];
  return (
    <span className={cn('rounded-md px-2 py-0.5 text-xs font-medium whitespace-nowrap', meta.className)}>
      {meta.label}
    </span>
  );
}
