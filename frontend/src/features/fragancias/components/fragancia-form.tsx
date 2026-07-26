'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { DecimalInput } from '@/components/forms/decimal-input';
import { FormField } from '@/components/forms/form-field';
import { SelectField } from '@/components/forms/select-field';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { getErrorMessage } from '@/lib/errors';
import { useCreateFragancia } from '../api/use-create-fragancia';
import { useUpdateFragancia } from '../api/use-update-fragancia';
import { fraganciaSchema, type FraganciaInput } from '../schemas/fragancia-schema';
import type { Fragancia } from '../types';

const GENERO_OPTIONS = [
  { value: 'femenina', label: 'Femenina' },
  { value: 'masculina', label: 'Masculina' },
];

export function FraganciaForm({ initialData }: { initialData?: Fragancia }) {
  const router = useRouter();
  const isEdit = Boolean(initialData);
  const create = useCreateFragancia();
  const update = useUpdateFragancia(initialData?.id ?? 0);

  const {
    register,
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FraganciaInput>({
    resolver: zodResolver(fraganciaSchema),
    defaultValues: {
      nombre_comercial: initialData?.nombre_comercial ?? '',
      nombre_alternativo: initialData?.nombre_alternativo ?? '',
      genero: initialData?.genero ?? 'femenina',
      gramos_minimo: initialData?.gramos_minimo ?? '',
    },
  });

  const onSubmit = async (data: FraganciaInput) => {
    const payload = {
      nombre_comercial: data.nombre_comercial,
      nombre_alternativo: data.nombre_alternativo ? data.nombre_alternativo : null,
      genero: data.genero,
      gramos_minimo: data.gramos_minimo,
    };
    try {
      if (isEdit) {
        await update.mutateAsync(payload);
        toast.success('Fragancia actualizada');
      } else {
        await create.mutateAsync(payload);
        toast.success('Fragancia creada');
      }
      router.push('/inventario/fragancias');
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Card>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
          <FormField id="nombre_comercial" label="Nombre comercial" error={errors.nombre_comercial?.message}>
            <Input id="nombre_comercial" disabled={isSubmitting} {...register('nombre_comercial')} />
          </FormField>

          <FormField
            id="nombre_alternativo"
            label="Nombre alternativo (opcional)"
            error={errors.nombre_alternativo?.message}
          >
            <Input id="nombre_alternativo" disabled={isSubmitting} {...register('nombre_alternativo')} />
          </FormField>

          <Controller
            control={control}
            name="genero"
            render={({ field }) => (
              <SelectField
                label="Género"
                value={field.value}
                onChange={field.onChange}
                options={GENERO_OPTIONS}
                error={errors.genero?.message}
                disabled={isSubmitting}
              />
            )}
          />

          <Controller
            control={control}
            name="gramos_minimo"
            render={({ field }) => (
              <FormField id="gramos_minimo" label="Gramos mínimo" error={errors.gramos_minimo?.message}>
                <DecimalInput
                  id="gramos_minimo"
                  suffix="g"
                  value={field.value}
                  onChange={field.onChange}
                  disabled={isSubmitting}
                />
              </FormField>
            )}
          />

          <div className="flex justify-end gap-2 pt-2">
            <Link href="/inventario/fragancias" className={buttonVariants({ variant: 'outline' })}>
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
