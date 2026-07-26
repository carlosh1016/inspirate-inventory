'use client';

import { Plus, X } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { formatCurrency, formatDateTime } from '@/lib/formatters';
import type { Consignacion, Cuadre, PagoCaja } from '../types';

interface Props {
  cuadre: Cuadre;
  /** admin+vendedora on an open cuadre may add pagos/consignaciones. */
  canManage: boolean;
  /** admin may delete pagos/consignaciones. */
  canDelete: boolean;
  onAddPago: () => void;
  onAddConsignacion: () => void;
  onDeletePago: (pago: PagoCaja) => void;
  onDeleteConsignacion: (consignacion: Consignacion) => void;
}

const METODOS: { label: string; key: keyof Cuadre }[] = [
  { label: 'Efectivo', key: 'total_efectivo' },
  { label: 'Nequi', key: 'total_nequi' },
  { label: 'Daviplata', key: 'total_daviplata' },
  { label: 'Transferencia', key: 'total_transferencia' },
  { label: 'Otros', key: 'total_otros' },
];

export function CuadreSections({
  cuadre,
  canManage,
  canDelete,
  onAddPago,
  onAddConsignacion,
  onDeletePago,
  onDeleteConsignacion,
}: Props) {
  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Ventas del día</CardTitle>
        </CardHeader>
        <CardContent>
          <dl className="space-y-1.5 text-sm">
            {METODOS.map((m) => (
              <div key={m.key} className="flex justify-between">
                <dt className="text-muted-foreground">{m.label}</dt>
                <dd className="tabular-nums">{formatCurrency(cuadre[m.key] as string)}</dd>
              </div>
            ))}
            <div className="flex justify-between border-t border-border pt-1.5 font-medium">
              <dt>Total ventas</dt>
              <dd className="tabular-nums">{formatCurrency(cuadre.total_ventas)}</dd>
            </div>
          </dl>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-row items-center justify-between">
          <CardTitle>Pagos del día</CardTitle>
          {canManage && (
            <Button variant="outline" size="sm" onClick={onAddPago}>
              <Plus className="size-4" />
              Registrar pago
            </Button>
          )}
        </CardHeader>
        <CardContent>
          {cuadre.pagos.length === 0 ? (
            <p className="text-sm text-muted-foreground">Sin pagos registrados.</p>
          ) : (
            <ul className="divide-y divide-border text-sm">
              {cuadre.pagos.map((pago) => (
                <li key={pago.id} className="flex items-center justify-between gap-3 py-2">
                  <span className="min-w-0">
                    <span className="font-medium">{pago.concepto}</span>
                    <span className="ml-2 text-muted-foreground">
                      {pago.usuario?.nombre_completo ?? ''} · {formatDateTime(pago.created_at)}
                    </span>
                  </span>
                  <span className="flex items-center gap-2">
                    <span className="tabular-nums">{formatCurrency(pago.monto)}</span>
                    {canDelete && (
                      <Button variant="ghost" size="icon-xs" aria-label="Eliminar pago" onClick={() => onDeletePago(pago)}>
                        <X className="size-3.5" />
                      </Button>
                    )}
                  </span>
                </li>
              ))}
            </ul>
          )}
          <div className="mt-2 flex justify-between border-t border-border pt-2 text-sm font-medium">
            <span>Total pagos</span>
            <span className="tabular-nums">{formatCurrency(cuadre.total_pagos)}</span>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex-row items-center justify-between">
          <CardTitle>Consignaciones</CardTitle>
          {canManage && (
            <Button variant="outline" size="sm" onClick={onAddConsignacion}>
              <Plus className="size-4" />
              Registrar consignación
            </Button>
          )}
        </CardHeader>
        <CardContent>
          {cuadre.consignaciones.length === 0 ? (
            <p className="text-sm text-muted-foreground">Sin consignaciones aún.</p>
          ) : (
            <ul className="divide-y divide-border text-sm">
              {cuadre.consignaciones.map((c) => (
                <li key={c.id} className="flex items-center justify-between gap-3 py-2">
                  <span className="min-w-0">
                    <span className="font-medium tabular-nums">{formatCurrency(c.monto)}</span>
                    <span className="ml-2 text-muted-foreground">
                      {[c.banco, c.referencia].filter(Boolean).join(' · ')}
                    </span>
                  </span>
                  {canDelete && (
                    <Button variant="ghost" size="icon-xs" aria-label="Eliminar consignación" onClick={() => onDeleteConsignacion(c)}>
                      <X className="size-3.5" />
                    </Button>
                  )}
                </li>
              ))}
            </ul>
          )}
          <div className="mt-2 flex justify-between border-t border-border pt-2 text-sm font-medium">
            <span>Total consignaciones</span>
            <span className="tabular-nums">{formatCurrency(cuadre.total_consignaciones)}</span>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="flex items-center justify-between py-4">
          <div>
            <p className="font-medium">Saldo esperado en caja</p>
            <p className="text-xs text-muted-foreground">
              Fondo base + efectivo − pagos − consignaciones
            </p>
          </div>
          <p className="text-2xl font-semibold tabular-nums">{formatCurrency(cuadre.saldo_calculado)}</p>
        </CardContent>
      </Card>
    </div>
  );
}
