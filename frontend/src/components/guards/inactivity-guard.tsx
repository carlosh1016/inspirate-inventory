'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { useLogout } from '@/features/auth/api/use-logout';
import { useInactivity } from '@/hooks/use-inactivity';
import { INACTIVITY_TIMEOUT_MS } from '@/lib/constants';
import { useAuthStore } from '@/stores/auth-store';

// Shows a modal after 1 hour of inactivity for a vendedora and logs her out on
// confirm. Admins have no inactivity timer.
export function InactivityGuard() {
  const router = useRouter();
  const usuario = useAuthStore((s) => s.usuario);
  const logout = useLogout();
  const [open, setOpen] = useState(false);

  const enabled = usuario?.rol === 'vendedora';
  useInactivity(INACTIVITY_TIMEOUT_MS, () => setOpen(true), enabled);

  const handleConfirm = async () => {
    await logout.mutateAsync();
    setOpen(false);
    router.push('/login');
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent showCloseButton={false}>
        <DialogHeader>
          <DialogTitle>Sesión expirada</DialogTitle>
          <DialogDescription>
            Tu sesión ha expirado por inactividad. Por seguridad, inicia sesión nuevamente.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button onClick={handleConfirm} disabled={logout.isPending}>
            Entendido
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
