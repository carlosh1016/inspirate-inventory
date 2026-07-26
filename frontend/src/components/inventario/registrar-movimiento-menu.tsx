'use client';

import { useRouter } from 'next/navigation';
import { ChevronDown, Plus } from 'lucide-react';

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { buttonVariants } from '@/components/ui/button';
import { usePermission } from '@/hooks/use-permission';

const BASE = '/inventario/movimientos';

export function RegistrarMovimientoMenu() {
  const router = useRouter();
  const { isAdmin } = usePermission();

  const items = [
    { label: 'Entrada de mercancía', href: `${BASE}/entrada`, adminOnly: false },
    { label: 'Traslado bodega → vitrina', href: `${BASE}/traslado`, adminOnly: false },
    { label: 'Producto dañado', href: `${BASE}/danado`, adminOnly: false },
    { label: 'Ajuste', href: `${BASE}/ajuste`, adminOnly: true },
    { label: 'Corrección', href: `${BASE}/correccion`, adminOnly: true },
  ].filter((item) => !item.adminOnly || isAdmin);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger className={buttonVariants()}>
        <Plus className="size-4" />
        Registrar movimiento
        <ChevronDown className="size-4" />
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {items.map((item) => (
          <DropdownMenuItem key={item.href} onClick={() => router.push(item.href)}>
            {item.label}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
