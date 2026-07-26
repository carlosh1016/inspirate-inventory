'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

import { cn } from '@/lib/utils';

const TABS = [
  { href: '/inventario/envases/modelos', label: 'Modelos' },
  { href: '/inventario/envases/variantes', label: 'Variantes' },
];

// Sub-navigation between modelos and variantes de envase.
export function EnvasesTabs() {
  const pathname = usePathname();
  return (
    <div className="mb-4 flex gap-1 border-b border-border">
      {TABS.map((tab) => {
        const active = pathname.startsWith(tab.href);
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
