'use client';

import { RequireRole } from '@/components/guards/require-role';
import { DashboardContent } from '@/features/dashboard/components/dashboard-content';

export default function DashboardPage() {
  return (
    <RequireRole role="admin">
      <DashboardContent />
    </RequireRole>
  );
}
