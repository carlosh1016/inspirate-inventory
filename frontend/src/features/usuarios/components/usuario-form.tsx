'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { FormField } from '@/components/forms/form-field';
import { SelectField } from '@/components/forms/select-field';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { getErrorMessage, parseApiError } from '@/lib/errors';
import type { UsuarioApi } from '@/types/domain';
import { useCreateUsuario } from '../api/use-create-usuario';
import { useUpdatePassword } from '../api/use-update-password';
import { useUpdateUsuario } from '../api/use-update-usuario';
import {
  createUsuarioSchema,
  editUsuarioSchema,
  type CreateUsuarioInput,
  type EditUsuarioInput,
} from '../schemas/usuario-schema';

const ROL_OPTIONS = [
  { value: 'admin', label: 'Admin' },
  { value: 'vendedora', label: 'Vendedora' },
];

const BASE_PATH = '/usuarios';

export function UsuarioForm({ initialData }: { initialData?: UsuarioApi }) {
  if (initialData) {
    return <EditForm usuario={initialData} />;
  }
  return <CreateForm />;
}

function CreateForm() {
  const router = useRouter();
  const create = useCreateUsuario();

  const {
    register,
    control,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<CreateUsuarioInput>({
    resolver: zodResolver(createUsuarioSchema),
    defaultValues: { nombre_completo: '', correo: '', password: '', confirmar_password: '', rol: 'vendedora' },
  });

  const onSubmit = async (data: CreateUsuarioInput) => {
    try {
      await create.mutateAsync(data);
      toast.success('Usuario creado');
      router.push(BASE_PATH);
    } catch (err) {
      const problem = parseApiError(err);
      if (problem.status === 409) {
        setError('correo', { message: 'Este email ya está registrado' });
        return;
      }
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Card>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
          <FormField id="nombre_completo" label="Nombre" error={errors.nombre_completo?.message}>
            <Input id="nombre_completo" disabled={isSubmitting} {...register('nombre_completo')} />
          </FormField>
          <FormField id="correo" label="Correo" error={errors.correo?.message}>
            <Input id="correo" type="email" disabled={isSubmitting} {...register('correo')} />
          </FormField>
          <FormField id="password" label="Contraseña" error={errors.password?.message}>
            <Input id="password" type="password" disabled={isSubmitting} {...register('password')} />
          </FormField>
          <FormField id="confirmar_password" label="Confirmar contraseña" error={errors.confirmar_password?.message}>
            <Input id="confirmar_password" type="password" disabled={isSubmitting} {...register('confirmar_password')} />
          </FormField>
          <Controller
            control={control}
            name="rol"
            render={({ field }) => (
              <SelectField
                label="Rol"
                value={field.value}
                onChange={field.onChange}
                options={ROL_OPTIONS}
                error={errors.rol?.message}
                disabled={isSubmitting}
              />
            )}
          />
          <div className="flex justify-end gap-2 pt-2">
            <Link href={BASE_PATH} className={buttonVariants({ variant: 'outline' })}>
              Cancelar
            </Link>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Guardando…' : 'Guardar'}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

function EditForm({ usuario }: { usuario: UsuarioApi }) {
  const router = useRouter();
  const update = useUpdateUsuario(usuario.id);
  const updatePassword = useUpdatePassword(usuario.id);

  const {
    register,
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<EditUsuarioInput>({
    resolver: zodResolver(editUsuarioSchema),
    defaultValues: {
      nombre_completo: usuario.nombre_completo,
      rol: usuario.rol,
      password_nueva: '',
      confirmar_password: '',
    },
  });

  const onSubmit = async (data: EditUsuarioInput) => {
    try {
      await update.mutateAsync({ nombre_completo: data.nombre_completo, rol: data.rol });
      if (data.password_nueva) {
        await updatePassword.mutateAsync(data.password_nueva);
      }
      toast.success('Usuario actualizado');
      router.push(BASE_PATH);
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Card>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
          <FormField id="nombre_completo" label="Nombre" error={errors.nombre_completo?.message}>
            <Input id="nombre_completo" disabled={isSubmitting} {...register('nombre_completo')} />
          </FormField>
          <FormField id="correo" label="Correo">
            <Input id="correo" value={usuario.correo} readOnly disabled className="cursor-not-allowed opacity-70" />
            <p className="text-xs text-muted-foreground">El correo no se puede cambiar.</p>
          </FormField>
          <Controller
            control={control}
            name="rol"
            render={({ field }) => (
              <SelectField
                label="Rol"
                value={field.value}
                onChange={field.onChange}
                options={ROL_OPTIONS}
                error={errors.rol?.message}
                disabled={isSubmitting}
              />
            )}
          />
          <div className="space-y-4 rounded-md border border-border p-4">
            <p className="text-sm font-medium">Cambiar contraseña (opcional)</p>
            <FormField id="password_nueva" label="Nueva contraseña" error={errors.password_nueva?.message}>
              <Input id="password_nueva" type="password" disabled={isSubmitting} {...register('password_nueva')} />
            </FormField>
            <FormField
              id="confirmar_password"
              label="Confirmar nueva contraseña"
              error={errors.confirmar_password?.message}
            >
              <Input id="confirmar_password" type="password" disabled={isSubmitting} {...register('confirmar_password')} />
            </FormField>
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Link href={BASE_PATH} className={buttonVariants({ variant: 'outline' })}>
              Cancelar
            </Link>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Guardando…' : 'Actualizar'}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
