'use client';

import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet';
import { useUIStore } from '@/stores/ui-store';
import { Sidebar } from './sidebar';

// The sidebar as a left drawer on screens < lg. Selecting an item closes it.
export function MobileDrawer() {
  const mobileDrawerOpen = useUIStore((s) => s.mobileDrawerOpen);
  const setMobileDrawerOpen = useUIStore((s) => s.setMobileDrawerOpen);

  return (
    <Sheet open={mobileDrawerOpen} onOpenChange={setMobileDrawerOpen}>
      <SheetContent side="left" className="w-[240px] p-0">
        <SheetHeader className="sr-only">
          <SheetTitle>Navegación</SheetTitle>
        </SheetHeader>
        <Sidebar onNavigate={() => setMobileDrawerOpen(false)} />
      </SheetContent>
    </Sheet>
  );
}
