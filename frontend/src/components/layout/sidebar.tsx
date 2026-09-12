'use client';

import { useState } from 'react';
import Image from 'next/image';
import { usePathname } from 'next/navigation';

import { useAuthStore } from '@/stores/auth-store';
import { navItems } from './nav-items';
import { SidebarItem } from './sidebar-item';

function isSectionActive(href: string, pathname: string): boolean {
  return href === '/dashboard' ? pathname === href : pathname === href || pathname.startsWith(`${href}/`);
}

function activeSectionHref(pathname: string): string | undefined {
  return navItems.find((item) => item.children && isSectionActive(item.href, pathname))?.href;
}

// Reusable sidebar content (240px). The desktop layout wraps it in a fixed
// <aside>; the mobile drawer renders it inside a Sheet (passing onNavigate to
// close on selection).
export function Sidebar({ onNavigate }: { onNavigate?: () => void }) {
  const usuario = useAuthStore((s) => s.usuario);
  const pathname = usePathname();

  // Which section a user has manually expanded/collapsed, overriding the
  // "expanded while its route is active" default. Reset when navigation
  // moves to a different section (or out of any section) — adjusting state
  // during render, per https://react.dev/learn/you-might-not-need-an-effect,
  // not in an effect, so this never trips a cascading re-render.
  const [manualOverride, setManualOverride] = useState<Map<string, boolean>>(new Map());
  const [prevPathname, setPrevPathname] = useState(pathname);
  if (pathname !== prevPathname) {
    if (activeSectionHref(pathname) !== activeSectionHref(prevPathname)) {
      setManualOverride(new Map());
    }
    setPrevPathname(pathname);
  }

  if (!usuario) return null;

  const visibleItems = navItems.filter((item) => item.roles.includes(usuario.rol));

  return (
    <div className="flex h-full w-[240px] flex-col bg-primary text-primary-foreground">
      <div className="border-b border-white/10 p-6">
        <div className="rounded-lg bg-white p-3">
          <Image
            src="/inspirate-logo.jpg"
            alt="Inspírate Perfumes & Cosmética"
            width={1542}
            height={688}
            priority
            className="h-auto w-full"
          />
        </div>
        <p className="mt-3 font-mono text-[10px] tracking-widest text-white/60 uppercase">
          Inventario
        </p>
      </div>
      <nav className="flex-1 space-y-1 overflow-y-auto px-3 py-4">
        {visibleItems.map((item) => {
          const expanded = manualOverride.has(item.href)
            ? manualOverride.get(item.href)!
            : isSectionActive(item.href, pathname);
          return (
            <SidebarItem
              key={item.href}
              item={item}
              active={isSectionActive(item.href, pathname)}
              pathname={pathname}
              rol={usuario.rol}
              highlight={usuario.rol === 'vendedora' && item.href === '/nueva-venta'}
              expanded={expanded}
              onToggleExpand={() =>
                setManualOverride((prev) => new Map(prev).set(item.href, !expanded))
              }
              onNavigate={onNavigate}
            />
          );
        })}
      </nav>
    </div>
  );
}
