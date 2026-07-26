'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

import { cn } from '@/lib/utils';

const TABS = [
  { href: '/inventario', label: 'Vista general', exact: true },
  { href: '/inventario/fragancias', label: 'Fragancias' },
  { href: '/inventario/envases', label: 'Envases' },
  { href: '/inventario/productos', label: 'Productos' },
  { href: '/inventario/movimientos', label: 'Movimientos' },
  { href: '/inventario/alertas', label: 'Alertas' },
];

export function InventarioTabs() {
  const pathname = usePathname();
  return (
    <div className="mb-4 flex flex-wrap gap-1 border-b border-border">
      {TABS.map((tab) => {
        const active = tab.exact ? pathname === tab.href : pathname.startsWith(tab.href);
        return (
          <Link
            key={tab.href}
            href={tab.href}
            className={cn(
              '-mb-px border-b-2 px-3 py-2 text-sm font-medium',
              active
                ? 'border-primary text-foreground'
                : 'border-transparent text-muted-foreground hover:text-foreground',
            )}
          >
            {tab.label}
          </Link>
        );
      })}
    </div>
  );
}
