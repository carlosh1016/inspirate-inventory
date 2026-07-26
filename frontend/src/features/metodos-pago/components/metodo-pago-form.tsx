'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { FormField } from '@/components/forms/form-field';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { getErrorMessage } from '@/lib/errors';
import { useCreateMetodoPago } from '../api/use-create-metodo-pago';
import { useUpdateMetodoPago } from '../api/use-update-metodo-pago';
import { metodoPagoSchema, type MetodoPagoInput } from '../schemas/metodo-pago-schema';
import type { MetodoPago } from '../types';

const BASE_PATH = '/configuracion/metodos-pago';

export function MetodoPagoForm({ initialData }: { initialData?: MetodoPago }) {
  const router = useRouter();
  const isEdit = Boolean(initialData);
  const create = useCreateMetodoPago();
  const update = useUpdateMetodoPago(initialData?.id ?? 0);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<MetodoPagoInput>({
    resolver: zodResolver(metodoPagoSchema),
    defaultValues: {
      nombre: initialData?.nombre ?? '',
      codigo: initialData?.codigo ?? '',
    },
  });

  const onSubmit = async (data: MetodoPagoInput) => {
    try {
      if (isEdit) {
        await update.mutateAsync(data);
        toast.success('Método de pago actualizado');
      } else {
        await create.mutateAsync(data);
        toast.success('Método de pago creado');
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
          <FormField id="nombre" label="Nombre" error={errors.nombre?.message}>
            <Input id="nombre" placeholder="Ej. Efectivo" disabled={isSubmitting} {...register('nombre')} />
          </FormField>
          <FormField id="codigo" label="Código" error={errors.codigo?.message}>
            <Input id="codigo" placeholder="Ej. efectivo" disabled={isSubmitting} {...register('codigo')} />
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
