'use client';

import { usePathname } from 'next/navigation';

import { useAuthStore } from '@/stores/auth-store';
import { navItems } from './nav-items';
import { SidebarItem } from './sidebar-item';

// Reusable sidebar content (240px). The desktop layout wraps it in a fixed
// <aside>; the mobile drawer renders it inside a Sheet (passing onNavigate to
// close on selection).
export function Sidebar({ onNavigate }: { onNavigate?: () => void }) {
  const usuario = useAuthStore((s) => s.usuario);
  const pathname = usePathname();

  if (!usuario) return null;

  const visibleItems = navItems.filter((item) => item.roles.includes(usuario.rol));

  return (
    <div className="flex h-full w-[240px] flex-col border-r border-border bg-background">
      <div className="border-b border-border p-6">
        <h1 className="text-xl font-semibold text-primary">Inspírate</h1>
        <p className="text-xs text-muted-foreground">Inventario</p>
      </div>
      <nav className="flex-1 space-y-1 overflow-y-auto px-3 py-4">
        {visibleItems.map((item) => (
          <SidebarItem
            key={item.href}
            item={item}
            active={
              item.href === '/dashboard'
                ? pathname === item.href
                : pathname === item.href || pathname.startsWith(`${item.href}/`)
            }
            highlight={usuario.rol === 'vendedora' && item.href === '/nueva-venta'}
            onNavigate={onNavigate}
          />
        ))}
      </nav>
    </div>
  );
}
