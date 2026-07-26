'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { FormField } from '@/components/forms/form-field';
import { MoneyInput } from '@/components/forms/money-input';
import { SelectField } from '@/components/forms/select-field';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { usePermission } from '@/hooks/use-permission';
import { getErrorMessage } from '@/lib/errors';
import { useCreateProducto } from '../api/use-create-producto';
import { useUpdateProducto } from '../api/use-update-producto';
import { productoSchema, type ProductoInput } from '../schemas/producto-schema';
import { CATEGORIA_OPTIONS, type Producto } from '../types';

const BASE_PATH = '/inventario/productos';

export function ProductoForm({ initialData }: { initialData?: Producto }) {
  const router = useRouter();
  const { isVendedora } = usePermission();
  const isEdit = Boolean(initialData);
  const create = useCreateProducto();
  const update = useUpdateProducto(initialData?.id ?? 0);

  // A vendedora may only change stock_minimo (backend-enforced); the rest are
  // shown disabled when she edits.
  const restricted = isEdit && isVendedora;

  const {
    register,
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ProductoInput>({
    resolver: zodResolver(productoSchema),
    defaultValues: {
      nombre: initialData?.nombre ?? '',
      categoria: initialData?.categoria ?? 'crema',
      precio: initialData?.precio ?? '',
      stock_minimo: initialData?.stock_minimo ?? 0,
    },
  });

  const onSubmit = async (data: ProductoInput) => {
    try {
      if (isEdit) {
        const payload = restricted
          ? { stock_minimo: data.stock_minimo }
          : {
              nombre: data.nombre,
              categoria: data.categoria,
              precio: data.precio,
              stock_minimo: data.stock_minimo,
            };
        await update.mutateAsync(payload);
        toast.success('Producto actualizado');
      } else {
        await create.mutateAsync({
          nombre: data.nombre,
          categoria: data.categoria,
          precio: data.precio,
          stock_minimo: data.stock_minimo,
        });
        toast.success('Producto creado');
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
          {restricted && (
            <p className="rounded-md bg-muted px-3 py-2 text-sm text-muted-foreground">
              Como vendedora solo puedes cambiar el stock mínimo.
            </p>
          )}

          <FormField id="nombre" label="Nombre" error={errors.nombre?.message}>
            <Input id="nombre" disabled={isSubmitting || restricted} {...register('nombre')} />
          </FormField>

          <Controller
            control={control}
            name="categoria"
            render={({ field }) => (
              <SelectField
                label="Categoría"
                value={field.value}
                onChange={field.onChange}
                options={CATEGORIA_OPTIONS}
                error={errors.categoria?.message}
                disabled={isSubmitting || restricted}
              />
            )}
          />

          <Controller
            control={control}
            name="precio"
            render={({ field }) => (
              <FormField id="precio" label="Precio" error={errors.precio?.message}>
                <MoneyInput
                  id="precio"
                  value={field.value}
                  onChange={field.onChange}
                  disabled={isSubmitting || restricted}
                />
              </FormField>
            )}
          />

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
