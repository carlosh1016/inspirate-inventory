'use client';

import { InactivityGuard } from '@/components/guards/inactivity-guard';
import { AuthProvider } from '@/components/layout/auth-provider';
import { MobileDrawer } from '@/components/layout/mobile-drawer';
import { Sidebar } from '@/components/layout/sidebar';
import { Topbar } from '@/components/layout/topbar';

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <div className="min-h-screen bg-background">
        <aside className="fixed top-0 left-0 hidden h-screen lg:block">
          <Sidebar />
        </aside>
        <MobileDrawer />
        <div className="lg:pl-[240px]">
          <Topbar />
          <main className="mx-auto max-w-7xl p-6">{children}</main>
        </div>
        <InactivityGuard />
      </div>
    </AuthProvider>
  );
}
