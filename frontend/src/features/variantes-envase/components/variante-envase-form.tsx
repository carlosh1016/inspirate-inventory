'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { Combobox } from '@/components/forms/combobox';
import { FormField } from '@/components/forms/form-field';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { getErrorMessage } from '@/lib/errors';
import { searchModelosEnvase } from '@/features/modelos-envase/api/search-modelos-envase';
import { useModelosLookup } from '@/features/modelos-envase/api/use-modelos-lookup';
import { useCreateVarianteEnvase } from '../api/use-create-variante-envase';
import { useUpdateVarianteEnvase } from '../api/use-update-variante-envase';
import { varianteEnvaseSchema, type VarianteEnvaseInput } from '../schemas/variante-envase-schema';
import type { VarianteEnvase } from '../types';

const BASE_PATH = '/inventario/envases/variantes';

export function VarianteEnvaseForm({ initialData }: { initialData?: VarianteEnvase }) {
  const router = useRouter();
  const isEdit = Boolean(initialData);
  const create = useCreateVarianteEnvase();
  const update = useUpdateVarianteEnvase(initialData?.id ?? 0);
  const { data: modelosMap } = useModelosLookup();
  const [modeloLabel, setModeloLabel] = useState<string | null>(null);

  const {
    register,
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<VarianteEnvaseInput>({
    resolver: zodResolver(varianteEnvaseSchema),
    defaultValues: {
      modelo_envase_id: initialData?.modelo_envase_id ?? 0,
      color: initialData?.color ?? '',
      stock_minimo: initialData?.stock_minimo ?? 0,
    },
  });

  const onSubmit = async (data: VarianteEnvaseInput) => {
    try {
      if (isEdit) {
        await update.mutateAsync({ color: data.color, stock_minimo: data.stock_minimo });
        toast.success('Variante actualizada');
      } else {
        await create.mutateAsync({
          modelo_envase_id: data.modelo_envase_id,
          color: data.color,
          stock_minimo: data.stock_minimo,
        });
        toast.success('Variante creada');
      }
      router.push(BASE_PATH);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  const modeloReadonlyLabel =
    initialData && (modelosMap?.get(initialData.modelo_envase_id) ?? `Modelo #${initialData.modelo_envase_id}`);

  return (
    <Card>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
          {isEdit ? (
            <div className="space-y-2">
              <Label>Modelo de envase</Label>
              <Input value={modeloReadonlyLabel ?? ''} disabled readOnly />
              <p className="text-xs text-muted-foreground">
                El modelo no se puede cambiar después de crear la variante.
              </p>
            </div>
          ) : (
            <Controller
              control={control}
              name="modelo_envase_id"
              render={({ field }) => (
                <FormField
                  id="modelo_envase_id"
                  label="Modelo de envase"
                  error={errors.modelo_envase_id?.message}
                >
                  <Combobox
                    value={field.value > 0 ? field.value : null}
                    selectedLabel={modeloLabel}
                    onChange={(id, label) => {
                      field.onChange(id ?? 0);
                      setModeloLabel(label);
                    }}
                    searchFn={searchModelosEnvase}
                    placeholder="Selecciona un modelo…"
                  />
                </FormField>
              )}
            />
          )}

          <FormField id="color" label="Color" error={errors.color?.message}>
            <Input id="color" disabled={isSubmitting} {...register('color')} />
          </FormField>

          <FormField id="stock_minimo" label="Stock mínimo" error={errors.stock_minimo?.message}>
            <Input
              id="stock_minimo"
              type="number"
              min={0}
              disabled={isSubmitting}
              {...register('stock_minimo', { valueAsNumber: true })}
            />
          </FormField>

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
