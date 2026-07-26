'use client';

import { useRouter } from 'next/navigation';
import { LogOut, User } from 'lucide-react';

import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useLogout } from '@/features/auth/api/use-logout';
import { useAuthStore } from '@/stores/auth-store';

function initialsOf(nombre: string): string {
  return nombre
    .trim()
    .split(/\s+/)
    .map((p) => p[0])
    .slice(0, 2)
    .join('')
    .toUpperCase();
}

export function UserMenu() {
  const router = useRouter();
  const usuario = useAuthStore((s) => s.usuario);
  const logout = useLogout();

  if (!usuario) return null;

  const handleLogout = async () => {
    await logout.mutateAsync();
    router.push('/login');
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        aria-label="Menú de usuario"
        className="flex items-center gap-2 rounded-md px-2 py-1 text-sm transition-colors hover:bg-muted focus-visible:ring-3 focus-visible:ring-ring/50 focus-visible:outline-none"
      >
        <Avatar className="size-8">
          <AvatarFallback className="text-xs">{initialsOf(usuario.nombre_completo)}</AvatarFallback>
        </Avatar>
        <span className="hidden md:inline">{usuario.nombre_completo}</span>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <div className="px-2 py-1.5">
          <p className="text-sm font-medium">{usuario.nombre_completo}</p>
          <p className="text-xs text-muted-foreground capitalize">{usuario.rol}</p>
        </div>
        <DropdownMenuSeparator />
        <DropdownMenuItem disabled>
          <User size={16} />
          Mi perfil
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleLogout}>
          <LogOut size={16} />
          Cerrar sesión
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
