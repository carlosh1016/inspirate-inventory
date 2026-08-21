'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { DecimalInput } from '@/components/forms/decimal-input';
import { FormField } from '@/components/forms/form-field';
import { MoneyInput } from '@/components/forms/money-input';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { getErrorMessage } from '@/lib/errors';
import { useCreateModeloEnvase } from '../api/use-create-modelo-envase';
import { useUpdateModeloEnvase } from '../api/use-update-modelo-envase';
import { modeloEnvaseSchema, type ModeloEnvaseInput } from '../schemas/modelo-envase-schema';
import type { ModeloEnvase } from '../types';

const BASE_PATH = '/inventario/envases/modelos';

export function ModeloEnvaseForm({ initialData }: { initialData?: ModeloEnvase }) {
  const router = useRouter();
  const isEdit = Boolean(initialData);
  const create = useCreateModeloEnvase();
  const update = useUpdateModeloEnvase(initialData?.id ?? 0);

  const {
    register,
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ModeloEnvaseInput>({
    resolver: zodResolver(modeloEnvaseSchema),
    defaultValues: {
      tipo: initialData?.tipo ?? '',
      tamano_oz: initialData?.tamano_oz ?? '',
      equiv_gramos: initialData?.equiv_gramos ?? '',
      precio_solo: initialData?.precio_solo ?? '',
      precio_con_fragancia: initialData?.precio_con_fragancia ?? '',
      precio_recarga: initialData?.precio_recarga ?? '',
      tiene_variantes: initialData?.tiene_variantes ?? true,
    },
  });

  const onSubmit = async (data: ModeloEnvaseInput) => {
    try {
      if (isEdit) {
        await update.mutateAsync({
          tipo: data.tipo,
          tamano_oz: data.tamano_oz,
          equiv_gramos: data.equiv_gramos,
          precio_solo: data.precio_solo,
          precio_con_fragancia: data.precio_con_fragancia,
          precio_recarga: data.precio_recarga,
        });
        toast.success('Modelo actualizado');
      } else {
        await create.mutateAsync(data);
        toast.success('Modelo creado');
      }
      router.push(BASE_PATH);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Card>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
          <FormField id="tipo" label="Tipo" error={errors.tipo?.message}>
            <Input id="tipo" placeholder="Ej. Spray, Roll-on" disabled={isSubmitting} {...register('tipo')} />
          </FormField>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Controller
              control={control}
              name="tamano_oz"
              render={({ field }) => (
                <FormField id="tamano_oz" label="Tamaño (oz)" error={errors.tamano_oz?.message}>
                  <DecimalInput id="tamano_oz" suffix="oz" value={field.value} onChange={field.onChange} disabled={isSubmitting} />
                </FormField>
              )}
            />
            <Controller
              control={control}
              name="equiv_gramos"
              render={({ field }) => (
                <FormField id="equiv_gramos" label="Equivalencia (gramos)" error={errors.equiv_gramos?.message}>
                  <DecimalInput id="equiv_gramos" suffix="g" value={field.value} onChange={field.onChange} disabled={isSubmitting} />
                </FormField>
              )}
            />
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <Controller
              control={control}
              name="precio_solo"
              render={({ field }) => (
                <FormField id="precio_solo" label="Precio solo" error={errors.precio_solo?.message}>
                  <MoneyInput id="precio_solo" value={field.value} onChange={field.onChange} disabled={isSubmitting} />
                </FormField>
              )}
            />
            <Controller
              control={control}
              name="precio_con_fragancia"
              render={({ field }) => (
                <FormField id="precio_con_fragancia" label="Con fragancia" error={errors.precio_con_fragancia?.message}>
                  <MoneyInput id="precio_con_fragancia" value={field.value} onChange={field.onChange} disabled={isSubmitting} />
                </FormField>
              )}
            />
            <Controller
              control={control}
              name="precio_recarga"
              render={({ field }) => (
                <FormField id="precio_recarga" label="Recarga" error={errors.precio_recarga?.message}>
                  <MoneyInput id="precio_recarga" value={field.value} onChange={field.onChange} disabled={isSubmitting} />
                </FormField>
              )}
            />
          </div>

          <Controller
            control={control}
            name="tiene_variantes"
            render={({ field }) => (
              <div className="flex items-start gap-2">
                <Checkbox
                  id="tiene_variantes"
                  checked={!field.value}
                  disabled={isSubmitting || isEdit}
                  onCheckedChange={(next) => field.onChange(next !== true)}
                />
                <div className="space-y-1">
                  <Label htmlFor="tiene_variantes" className="font-normal">
                    Este modelo no tiene variantes (ej. envase de lujo)
                  </Label>
                  <p className="text-xs text-muted-foreground">
                    Márcalo solo si este envase es único y no varía por grosor. No se puede cambiar después de
                    creado.
                  </p>
                </div>
              </div>
            )}
          />

          <div className="flex justify-end gap-2 pt-2">
            <Link href={BASE_PATH} className={buttonVariants({ variant: 'outline' })}>
              Cancelar
            </Link>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Guardando…' : isEdit ? 'Actualizar' : 'Guardar'}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
