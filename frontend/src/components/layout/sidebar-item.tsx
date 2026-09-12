'use client';

import Link from 'next/link';
import { ChevronDown } from 'lucide-react';

import { cn } from '@/lib/utils';
import type { Rol } from '@/types/domain';
import type { NavItem } from './nav-items';

interface Props {
  item: NavItem;
  active: boolean;
  pathname: string;
  rol: Rol;
  /** Emphasize this item (e.g. "Nueva venta" for a vendedora). */
  highlight?: boolean;
  expanded?: boolean;
  onToggleExpand?: () => void;
  onNavigate?: () => void;
}

function isSubItemActive(href: string, pathname: string): boolean {
  return href === '/inventario' ? pathname === href : pathname === href || pathname.startsWith(`${href}/`);
}

export function SidebarItem({
  item,
  active,
  pathname,
  rol,
  highlight = false,
  expanded = false,
  onToggleExpand,
  onNavigate,
}: Props) {
  const Icon = item.icon;
  const children = item.children?.filter((c) => !c.roles || c.roles.includes(rol));
  const hasChildren = !!children && children.length > 0;

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
      {hasChildren && !item.disabled && (
        <button
          type="button"
          onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            onToggleExpand?.();
          }}
          aria-label={expanded ? `Contraer ${item.label}` : `Expandir ${item.label}`}
          aria-expanded={expanded}
          className="-my-1 -mr-1 ml-auto rounded p-1 hover:bg-white/10"
        >
          <ChevronDown size={14} className={cn('transition-transform', expanded && 'rotate-180')} />
        </button>
      )}
    </div>
  );

  const parent = item.disabled ? (
    <div aria-disabled="true">{content}</div>
  ) : (
    <Link href={item.href} onClick={onNavigate} aria-current={active ? 'page' : undefined}>
      {content}
    </Link>
  );

  if (!hasChildren) return parent;

  return (
    <div>
      {parent}
      {expanded && (
        <div className="mt-1 ml-4 space-y-0.5 border-l border-white/10 pl-4">
          {children.map((sub) => {
            const subActive = isSubItemActive(sub.href, pathname);
            return (
              <Link
                key={sub.href}
                href={sub.href}
                onClick={onNavigate}
                aria-current={subActive ? 'page' : undefined}
                className={cn(
                  'block rounded-md px-3 py-1.5 text-sm transition-colors',
                  subActive ? 'font-semibold text-highlight' : 'text-white/60 hover:bg-white/5 hover:text-white',
                )}
              >
                {sub.label}
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}
