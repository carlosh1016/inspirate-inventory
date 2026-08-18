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
        'flex items-center gap-3 rounded-md border-r-2 border-transparent px-3 py-2 text-sm transition-colors',
        active && 'bg-white/10 font-semibold text-highlight',
        !active && highlight && !item.disabled && 'bg-highlight font-semibold text-highlight-foreground hover:bg-highlight/90',
        !active && !highlight && !item.disabled && 'text-white/70 hover:bg-white/5 hover:text-white',
        item.disabled && 'cursor-not-allowed text-white/30',
      )}
    >
      <Icon size={18} />
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
