'use client';

import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';

import { FormField } from '@/components/forms/form-field';
import { PasswordInput } from '@/components/forms/password-input';
import { Button, buttonVariants } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { getErrorMessage } from '@/lib/errors';
import { cn } from '@/lib/utils';
import { useConfirmPasswordReset } from '../api/use-confirm-password-reset';
import { resetPasswordSchema, type ResetPasswordInput } from '../schemas/reset-password-schema';

export function ResetPasswordForm() {
  const router = useRouter();
  const params = useSearchParams();
  const token = params.get('token') ?? '';
  const confirm = useConfirmPasswordReset();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ResetPasswordInput>({ resolver: zodResolver(resetPasswordSchema) });

  if (!token) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Enlace inválido</CardTitle>
          <CardDescription>El enlace de restablecimiento no es válido o está incompleto.</CardDescription>
        </CardHeader>
        <CardContent>
          <Link href="/forgot-password" className={cn(buttonVariants(), 'w-full')}>
            Solicitar un nuevo enlace
          </Link>
        </CardContent>
      </Card>
    );
  }

  const onSubmit = async (data: ResetPasswordInput) => {
    try {
      await confirm.mutateAsync({ token, password_nueva: data.password_nueva });
      toast.success('Contraseña actualizada. Ya puedes iniciar sesión.');
      router.push('/login');
    } catch (err) {
      toast.error(getErrorMessage(err));
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Nueva contraseña</CardTitle>
        <CardDescription>Crea una contraseña nueva para tu cuenta.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4" noValidate>
          <FormField id="password_nueva" label="Nueva contraseña" error={errors.password_nueva?.message}>
            <PasswordInput
              id="password_nueva"
              autoComplete="new-password"
              disabled={isSubmitting}
              {...register('password_nueva')}
            />
          </FormField>

          <FormField
            id="confirmar_password"
            label="Confirmar contraseña"
            error={errors.confirmar_password?.message}
          >
            <PasswordInput
              id="confirmar_password"
              autoComplete="new-password"
              disabled={isSubmitting}
              {...register('confirmar_password')}
            />
          </FormField>

          <Button type="submit" className="w-full" disabled={isSubmitting}>
            {isSubmitting ? 'Guardando...' : 'Cambiar contraseña'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
