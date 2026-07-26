'use client';

import Link from 'next/link';
import { AlertTriangle } from 'lucide-react';

import { usePermission } from '@/hooks/use-permission';
import { formatDateShort } from '@/lib/formatters';
import { useCuadreAnteriorAbierto } from '../api/use-cuadre-anterior-abierto';

// Persistent, reload-safe banner (admin only) when a previous day's cuadre was
// left open. Detected from the cuadres history, so it survives reloads.
export function CuadreWarningAnterior() {
  const { isAdmin } = usePermission();
  const { data: anterior } = useCuadreAnteriorAbierto(isAdmin);

  if (!anterior) return null;

  return (
    <div className="mb-4 flex flex-wrap items-center justify-between gap-2 rounded-md border border-warning/30 bg-warning/10 px-3 py-2 text-sm text-warning">
      <span className="flex items-center gap-2">
        <AlertTriangle className="size-4 shrink-0" />
        El cuadre del día {formatDateShort(anterior.fecha)} quedó abierto sin cerrar.
      </span>
      <Link href={`/caja/${anterior.id}`} className="font-medium underline underline-offset-2">
        Ciérralo ahora →
      </Link>
    </div>
  );
}
