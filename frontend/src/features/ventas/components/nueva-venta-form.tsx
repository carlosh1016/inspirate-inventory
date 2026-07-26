'use client';

import { useMemo, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { toast } from 'sonner';

import { FormError } from '@/components/forms/form-error';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Textarea } from '@/components/ui/textarea';
import { ErrorState } from '@/components/feedback/error-state';
import { useModelosFullLookup } from '@/features/modelos-envase/api/use-modelos-full-lookup';
import { useCreateVenta, type CreateVentaItemPayload, type CreateVentaPayload } from '../api/use-create-venta';
import { useMetodosPagoActivos } from '../api/use-metodos-pago-activos';
import { useCalculoVenta, itemSubtotal } from '../hooks/use-calculo-venta';
import { useIdempotencyKey } from '../hooks/use-idempotency-key';
import { nuevaVentaSchema } from '../schemas/nueva-venta-schema';
import { parseVentaError } from '../venta-errors';
import type { TipoLinea, VentaItemState } from '../types';
import { VentaConfirmarDialog } from './venta-confirmar-dialog';
import { VentaItemRow } from './venta-item-row';
import { VentaTipoLineaSelector } from './venta-tipo-linea-selector';
import { VentaTotales } from './venta-totales';

function emptyItem(tipo: TipoLinea): VentaItemState {
  return {
    key: crypto.randomUUID(),
    tipo_linea: tipo,
    envase: null,
    fragancia: null,
    producto: null,
    gramos: '',
    cantidad: 1,
    feromona_enabled: false,
    feromona: null,
  };
}

function esFragancia(tipo: TipoLinea) {
  return tipo === 'envase_con_fragancia' || tipo === 'recarga';
}

function toItemPayload(it: VentaItemState): CreateVentaItemPayload {
  switch (it.tipo_linea) {
    case 'envase_con_fragancia':
    case 'recarga':
      return {
        tipo_linea: it.tipo_linea,
        variante_envase_id: it.envase?.variante_id ?? 0,
        fragancia_id: it.fragancia?.id ?? 0,
        gramos_fragancia: it.gramos,
        cantidad: it.cantidad,
        feromona_producto_id: it.feromona_enabled && it.feromona ? it.feromona.id : undefined,
      };
    case 'envase_solo':
      return { tipo_linea: 'envase_solo', variante_envase_id: it.envase?.variante_id ?? 0, cantidad: it.cantidad };
    case 'producto_otro':
      return { tipo_linea: 'producto_otro', producto_id: it.producto?.id ?? 0, cantidad: it.cantidad };
  }
}

export function NuevaVentaForm() {
  const router = useRouter();
  const idempotencyKey = useIdempotencyKey();
  const create = useCreateVenta();
  const { data: modelosMap, isError: modelosError } = useModelosFullLookup();
  const { data: metodos, isError: metodosError } = useMetodosPagoActivos();

  const [items, setItems] = useState<VentaItemState[]>([]);
  const [metodoPagoId, setMetodoPagoId] = useState(0);
  const [observaciones, setObservaciones] = useState('');
  const [errorByIndex, setErrorByIndex] = useState<Map<number, string>>(new Map());
  const [metodoError, setMetodoError] = useState<string | undefined>();
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [cuadreCerrado, setCuadreCerrado] = useState(false);

  const calculo = useCalculoVenta(items);
  const subtotals = useMemo(() => items.map(itemSubtotal), [items]);
  const metodoNombre = metodos?.find((m) => m.id === metodoPagoId)?.nombre ?? '';

  const addItem = (tipo: TipoLinea) => setItems((prev) => [...prev, emptyItem(tipo)]);
  const updateItem = (index: number, next: VentaItemState) =>
    setItems((prev) => prev.map((it, i) => (i === index ? next : it)));
  const removeItem = (index: number) => setItems((prev) => prev.filter((_, i) => i !== index));

  const validate = () => {
    const errs = new Map<number, string>();
    items.forEach((it, i) => {
      if (esFragancia(it.tipo_linea) && it.feromona_enabled && !it.feromona) {
        errs.set(i, 'Selecciona la feromona o desactiva el checkbox');
      }
    });
    const payload: CreateVentaPayload = {
      metodo_pago_id: metodoPagoId,
      observaciones: observaciones.trim() ? observaciones.trim() : null,
      items: items.map(toItemPayload),
    };
    const parsed = nuevaVentaSchema.safeParse(payload);
    let metodoErr: string | undefined;
    if (!parsed.success) {
      for (const issue of parsed.error.issues) {
        const [p0, p1] = issue.path;
        if (p0 === 'items' && typeof p1 === 'number') {
          if (!errs.has(p1)) errs.set(p1, issue.message);
        } else if (p0 === 'metodo_pago_id') {
          metodoErr = issue.message;
        }
      }
    }
    return { errs, metodoErr, payload, ok: parsed.success && errs.size === 0 };
  };

  const submit = async (payload: CreateVentaPayload) => {
    try {
      const venta = await create.mutateAsync({ payload, idempotencyKey });
      toast.success(`Venta ${venta.consecutivo} registrada correctamente`);
      router.push(`/ventas/${venta.id}`);
    } catch (err) {
      const info = parseVentaError(err);
      setConfirmOpen(false);
      if (info.kind === 'stock' || info.kind === 'coherence') {
        setErrorByIndex(info.byIndex);
        toast.error(info.message);
      } else if (info.kind === 'cuadre_cerrado') {
        setCuadreCerrado(true);
        toast.error(info.message);
      } else {
        toast.error(info.message);
      }
    }
  };

  const handleRegistrar = () => {
    const result = validate();
    setErrorByIndex(result.errs);
    setMetodoError(result.metodoErr);
    if (items.length === 0) {
      toast.error('Agrega al menos un ítem');
      return;
    }
    if (!result.ok) return;
    if (calculo.total > 200000) {
      setConfirmOpen(true);
      return;
    }
    void submit(result.payload);
  };

  if (modelosError || metodosError) {
    return <ErrorState error={new Error('No se pudieron cargar los datos del catálogo.')} />;
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Ítems</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <VentaTipoLineaSelector onAdd={addItem} />
          {items.length === 0 ? (
            <p className="rounded-lg border border-dashed border-border py-8 text-center text-sm text-muted-foreground">
              Agrega un ítem con los botones de arriba para empezar.
            </p>
          ) : (
            items.map((item, index) => (
              <VentaItemRow
                key={item.key}
                index={index}
                value={item}
                onChange={(next) => updateItem(index, next)}
                onRemove={() => removeItem(index)}
                modelosMap={modelosMap ?? new Map()}
                subtotal={subtotals[index] ?? 0}
                error={errorByIndex.get(index)}
              />
            ))
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Método de pago</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <RadioGroup
            value={metodoPagoId > 0 ? String(metodoPagoId) : ''}
            onValueChange={(v) => setMetodoPagoId(Number(v))}
            className="flex flex-row flex-wrap gap-4"
          >
            {(metodos ?? []).map((m) => (
              <Label key={m.id} className="cursor-pointer gap-2 font-normal">
                <RadioGroupItem value={String(m.id)} />
                {m.nombre}
              </Label>
            ))}
          </RadioGroup>
          <FormError message={metodoError} />
        </CardContent>
      </Card>

      <Card>
        <CardContent className="space-y-4">
          <VentaTotales calculo={calculo} />
          <div className="space-y-2">
            <Label htmlFor="observaciones">Observaciones (opcional)</Label>
            <Textarea
              id="observaciones"
              value={observaciones}
              onChange={(e) => setObservaciones(e.target.value)}
              maxLength={1000}
            />
          </div>
          {cuadreCerrado && (
            <p className="rounded-md bg-destructive/10 px-3 py-2 text-sm text-destructive">
              El cuadre de caja del día está cerrado. No se pueden registrar más ventas. Vuelve mañana o
              pide al admin que reabra el cuadre.
            </p>
          )}
          <div className="flex justify-end gap-2">
            <Link href="/ventas" className={buttonVariants({ variant: 'outline' })}>
              Cancelar
            </Link>
            <Button onClick={handleRegistrar} disabled={create.isPending || cuadreCerrado}>
              {create.isPending ? 'Registrando…' : 'Registrar venta'}
            </Button>
          </div>
        </CardContent>
      </Card>

      <VentaConfirmarDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        total={calculo.total}
        itemsCount={items.length}
        metodoPagoNombre={metodoNombre}
        submitting={create.isPending}
        onConfirm={() => {
          const payload: CreateVentaPayload = {
            metodo_pago_id: metodoPagoId,
            observaciones: observaciones.trim() ? observaciones.trim() : null,
            items: items.map(toItemPayload),
          };
          void submit(payload);
        }}
      />
    </div>
  );
}
