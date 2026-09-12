import {
  Clock,
  FileText,
  LayoutDashboard,
  Package,
  Settings,
  ShieldCheck,
  ShoppingCart,
  Users,
  Wallet,
  type LucideIcon,
} from 'lucide-react';

import type { Rol } from '@/types/domain';

export interface NavSubItem {
  href: string;
  label: string;
  /** Defaults to the parent's `roles` when omitted. */
  roles?: Rol[];
}

export interface NavItem {
  href: string;
  label: string;
  icon: LucideIcon;
  roles: Rol[];
  /** Rutas de tandas futuras: se muestran con badge "próximamente" y sin link. */
  disabled?: boolean;
  /** Secciones internas mostradas como submenú desplegable bajo este ítem. */
  children?: NavSubItem[];
}

const INVENTARIO_CHILDREN: NavSubItem[] = [
  { href: '/inventario', label: 'Vista general' },
  { href: '/inventario/fragancias', label: 'Fragancias' },
  { href: '/inventario/envases', label: 'Envases' },
  { href: '/inventario/productos', label: 'Productos' },
  { href: '/inventario/movimientos', label: 'Movimientos' },
  { href: '/inventario/conteo-mensual', label: 'Conteo mensual', roles: ['admin'] },
  { href: '/inventario/alertas', label: 'Alertas' },
];

// Only Dashboard is live this tanda; everything else is a disabled placeholder.
export const navItems: NavItem[] = [
  { href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard, roles: ['admin'] },
  { href: '/nueva-venta', label: 'Nueva venta', icon: ShoppingCart, roles: ['admin', 'vendedora'] },
  { href: '/mi-turno', label: 'Mi turno', icon: Clock, roles: ['vendedora'] },
  {
    href: '/inventario',
    label: 'Inventario',
    icon: Package,
    roles: ['admin', 'vendedora'],
    children: INVENTARIO_CHILDREN,
  },
  { href: '/ventas', label: 'Ventas', icon: ShoppingCart, roles: ['admin', 'vendedora'] },
  { href: '/caja', label: 'Caja', icon: Wallet, roles: ['admin', 'vendedora'] },
  { href: '/usuarios', label: 'Usuarios', icon: Users, roles: ['admin'] },
  { href: '/reportes', label: 'Reportes', icon: FileText, roles: ['admin'], disabled: true },
  { href: '/auditoria', label: 'Auditoría', icon: ShieldCheck, roles: ['admin'], disabled: true },
  { href: '/configuracion', label: 'Configuración', icon: Settings, roles: ['admin'] },
];
