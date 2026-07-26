import { Banknote, CreditCard, Landmark, Smartphone, type LucideIcon } from 'lucide-react';

import { cn } from '@/lib/utils';

const ICON_BY_CODIGO: Record<string, LucideIcon> = {
  efectivo: Banknote,
  nequi: Smartphone,
  daviplata: Smartphone,
  transferencia: Landmark,
};

export function MetodoPagoBadge({
  nombre,
  codigo,
  className,
}: {
  nombre: string;
  codigo?: string;
  className?: string;
}) {
  const Icon = (codigo && ICON_BY_CODIGO[codigo]) || CreditCard;
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 rounded-md bg-muted px-2 py-0.5 text-xs font-medium',
        className,
      )}
    >
      <Icon className="size-3.5 text-muted-foreground" />
      {nombre}
    </span>
  );
}
