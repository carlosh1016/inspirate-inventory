'use client';

import { Menu } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { useUIStore } from '@/stores/ui-store';
import { UserMenu } from './user-menu';

export function Topbar() {
  const setMobileDrawerOpen = useUIStore((s) => s.setMobileDrawerOpen);

  return (
    <header className="sticky top-0 z-10 flex h-14 items-center justify-between border-b border-border bg-background px-4 lg:px-6">
      <Button
        variant="ghost"
        size="icon"
        className="lg:hidden"
        aria-label="Abrir menú"
        onClick={() => setMobileDrawerOpen(true)}
      >
        <Menu size={20} />
      </Button>
      <div className="flex-1" />
      <UserMenu />
    </header>
  );
}
