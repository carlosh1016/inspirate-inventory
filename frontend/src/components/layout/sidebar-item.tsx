'use client';

import Link from 'next/link';

import { cn } from '@/lib/utils';
import type { NavItem } from './nav-items';

interface Props {
  item: NavItem;
  active: boolean;
  /** Emphasize this item (e.g. "Nueva venta" for a vendedora). */
  highlight?: boolean;
  onNavigate?: () => void;
}

export function SidebarItem({ item, active, highlight = false, onNavigate }: Props) {
  const Icon = item.icon;

  const content = (
    <div
      className={cn(
        'flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors',
        active && 'bg-muted font-medium text-foreground',
        !active && highlight && !item.disabled && 'bg-primary/10 font-medium text-primary hover:bg-primary/15',
        !active && !highlight && !item.disabled && 'text-muted-foreground hover:bg-muted hover:text-foreground',
        item.disabled && 'cursor-not-allowed text-muted-foreground/50',
      )}
    >
      <Icon size={18} className={cn(!active && highlight && 'text-primary')} />
      <span>{item.label}</span>
      {item.disabled && (
        <span className="ml-auto text-[10px] tracking-wider uppercase">próximamente</span>
      )}
    </div>
  );

  if (item.disabled) {
    return <div aria-disabled="true">{content}</div>;
  }

  return (
    <Link href={item.href} onClick={onNavigate} aria-current={active ? 'page' : undefined}>
      {content}
    </Link>
  );
}
